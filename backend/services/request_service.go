package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
)

type RequestService struct {
	requestRepo         repositories.IRequestRepo
	templateRepo        repositories.IRequestTemplateRepo
	expectedDocRepo     repositories.IExpectedDocumentRepo
	expectedDocTmplRepo repositories.IExpectedDocumentTemplateRepo
	txManager           repositories.ITxManager
	fileStorage         IFileStorageService
	logger              *slog.Logger
}

func NewRequestService(
	requestRepo repositories.IRequestRepo,
	templateRepo repositories.IRequestTemplateRepo,
	expectedDocRepo repositories.IExpectedDocumentRepo,
	expectedDocTmplRepo repositories.IExpectedDocumentTemplateRepo,
	txManager repositories.ITxManager,
	logger *slog.Logger,
	fileStorage IFileStorageService,
) *RequestService {
	return &RequestService{
		requestRepo:         requestRepo,
		templateRepo:        templateRepo,
		expectedDocRepo:     expectedDocRepo,
		expectedDocTmplRepo: expectedDocTmplRepo,
		txManager:           txManager,
		logger:              logger,
		fileStorage:         fileStorage,
	}
}

func (service *RequestService) AddRequest(ctx context.Context, claims types.JWTClaims, dto models.RequestDTOCreate) (*int, error) {
	if claims.IsAdmin() || claims.IsDepartmentMember() {
		service.logger.Warn("admin or department member tried to create request for themselves, got rejected",
			slog.Int("jwt_user_id", claims.UserID),
		)
		return nil, errors.ErrForbidden{Msg: "You are not allowed to create requests"}
	}

	if err := ValidateRequestInput(dto); err != nil {
		service.logger.Warn("request creation failed because it did not pass validations",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	template, err := service.templateRepo.GetRequestTemplateByID(ctx, dto.TemplateID)
	if err != nil {
		service.logger.Warn("template not found for request creation",
			slog.Int("template_id", dto.TemplateID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, errors.ErrNotFound{Msg: "Template not found."}
	}

	expectedDocTemplates, err := service.expectedDocTmplRepo.GetByRequestTemplateID(ctx, dto.TemplateID)
	if err != nil {
		service.logger.Error("failed to fetch expected document templates",
			slog.Int("template_id", dto.TemplateID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	expectedDocs := make([]models.ExpectedDocument, len(expectedDocTemplates))
	for i, edt := range expectedDocTemplates {
		expectedDocs[i] = models.ExpectedDocument{
			Title:           edt.Title,
			Description:     edt.Description,
			Status:          "pending",
			ExampleFilePath: edt.ExampleFilePath,
			ExampleMimeType: edt.ExampleMimeType,
		}
	}

	nextDueAt := ComputeNextDueAt(dto.DueDate, template.RecurrenceCron)
	req := models.Request{
		RequestBase: models.RequestBase{
			Assignee:       claims.UserID,
			DepartmentID:   template.DepartmentID,
			Title:          template.Title,
			Description:    template.Description,
			IsRecurring:    template.IsRecurring,
			RecurrenceCron: template.RecurrenceCron,
			IsScheduled:    dto.IsScheduled,
			ScheduledFor:   dto.ScheduledFor,
			NextDueAt:      nextDueAt,
			LastUploadedAt: nil,
			DueDate:        dto.DueDate,
		},
		RequestTemplateID: &dto.TemplateID,
	}

	id, err := service.createRequestTransaction(ctx, req, expectedDocs)
	if err != nil {
		service.logger.Error("failed to create request",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	service.logger.Info("request created successfully",
		slog.Int("request_id", *id),
		slog.Int("jwt_user_id", claims.UserID),
		slog.Int("template_id", dto.TemplateID),
		slog.Int("department_id", template.DepartmentID),
	)
	return id, nil
}

func (service *RequestService) GetAllRequests(ctx context.Context, claims types.JWTClaims, search *string) ([]models.RequestDTORead, error) {
	if !claims.IsAdmin() {
		service.logger.Warn("unauthorized attempt to get all requests",
			slog.Int("jwt_user_id", claims.UserID),
		)
		return nil, errors.ErrForbidden{Msg: "Only admins can view all requests."}
	}

	reqs, err := service.requestRepo.GetAllRequests(ctx, search)
	if err != nil {
		service.logger.Error("failed to get all requests",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	for i := range reqs {
		reqs[i].Status = ComputeStatus(reqs[i].LastUploadedAt, reqs[i].NextDueAt, reqs[i].ExpectedDocuments)
	}

	service.logger.Info("fetched all requests successfully",
		slog.Int("jwt_user_id", claims.UserID),
	)
	return reqs, nil
}

func (s *RequestService) GetRequestByID(ctx context.Context, claims types.JWTClaims, requestID int) (*models.RequestDTORead, error) {
	req, err := s.checkUserIsParticipantOfRequest(ctx, claims, requestID)
	if err != nil {
		return nil, err
	}

	expectedDocs, err := s.expectedDocRepo.GetExpectedDocumentsByRequest(ctx, requestID)
	if err != nil {
		s.logger.Error("failed to get expected documents for request",
			slog.Int("request_id", requestID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}
	req.ExpectedDocuments = expectedDocs
	req.Status = ComputeStatus(req.LastUploadedAt, req.NextDueAt, req.ExpectedDocuments)

	return req, nil
}

func (service *RequestService) GetRequestsByAssignee(ctx context.Context, claims types.JWTClaims, assigneeID int, search *string) ([]models.RequestDTORead, error) {
	if !claims.IsAdmin() && assigneeID != claims.UserID {
		service.logger.Warn("unauthorized access attempt to requests by assignee",
			slog.Int("assignee_id", assigneeID),
			slog.Int("jwt_user_id", claims.UserID),
		)
		return nil, errors.ErrForbidden{Msg: "You are not allowed to view these requests."}
	}

	reqs, err := service.requestRepo.GetRequestsByAssigneeWithExpectedDocs(ctx, assigneeID, search)
	if err != nil {
		service.logger.Error("error when retrieving requests by assignee",
			slog.Int("assignee_id", assigneeID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	for i := range reqs {
		reqs[i].Status = ComputeStatus(reqs[i].LastUploadedAt, reqs[i].NextDueAt, reqs[i].ExpectedDocuments)
	}

	service.logger.Info("retrieved requests by assignee successfully",
		slog.Int("assignee_id", assigneeID),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return reqs, nil
}

func (service *RequestService) GetRequestsByDepartment(ctx context.Context, claims types.JWTClaims, departmentID int, search *string) ([]models.RequestDTORead, error) {
	isMemberOfDepartment := claims.DepartmentID != nil && *claims.DepartmentID == departmentID
	if !claims.IsAdmin() && !isMemberOfDepartment {
		service.logger.Warn("unauthorized access attempt to requests by department",
			slog.Int("department_id", departmentID),
			slog.Int("jwt_user_id", claims.UserID),
		)
		return nil, errors.ErrForbidden{Msg: "You are not allowed to view these requests."}
	}

	reqs, err := service.requestRepo.GetRequestsByDepartmentWithExpectedDocs(ctx, departmentID, search)
	if err != nil {
		service.logger.Error("error when retrieving requests by department",
			slog.Int("department_id", departmentID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	for i := range reqs {
		reqs[i].Status = ComputeStatus(reqs[i].LastUploadedAt, reqs[i].NextDueAt, reqs[i].ExpectedDocuments)
	}

	service.logger.Info("retrieved requests by department successfully",
		slog.Int("department_id", departmentID),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return reqs, nil
}

func (service *RequestService) ForwardRequestToDepartment(ctx context.Context, claims types.JWTClaims, requestID int, departmentID int) error {
	if !claims.IsAdmin() {
		service.logger.Warn("unauthorized attempt to forward request",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("request_id", requestID),
		)
		return errors.ErrForbidden{Msg: "Only admins can forward requests to departments."}
	}

	if err := service.requestRepo.ForwardRequestToDepartment(ctx, requestID, departmentID); err != nil {
		service.logger.Error("failed to forward request to department",
			slog.Int("request_id", requestID),
			slog.Int("department_id", departmentID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return err
	}

	service.logger.Info("request forwarded to department successfully",
		slog.Int("request_id", requestID),
		slog.Int("department_id", departmentID),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return nil
}

func (s *RequestService) PatchRequest(ctx context.Context, claims types.JWTClaims, requestID int, updatedDTO models.RequestDTOPatch) error {
	if err := ValidatePatchDTO(updatedDTO); err != nil {
		s.logger.Warn("request patch validation failed",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("request_id", requestID),
			slog.Any("error", err),
		)
		return err
	}

	if _, err := s.checkUserCanEditRequest(ctx, claims, requestID); err != nil {
		return err
	}

	if err := s.requestRepo.UpdateRequestTitle(ctx, requestID, updatedDTO.Title); err != nil {
		s.logger.Error("failed to update request title",
			slog.Int("request_id", requestID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.String("new_title", updatedDTO.Title),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("request patched successfully",
		slog.Int("request_id", requestID),
		slog.Int("jwt_user_id", claims.UserID),
		slog.String("new_title", updatedDTO.Title),
	)
	return nil
}

func (s *RequestService) ReopenRequest(ctx context.Context, claims types.JWTClaims, requestID int) error {
	if _, err := s.checkUserCanEditRequest(ctx, claims, requestID); err != nil {
		return err
	}

	if err := s.requestRepo.ReopenRequest(ctx, requestID); err != nil {
		s.logger.Error("error when trying to reopen request",
			slog.Int("request_id", requestID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("reopened request successfully",
		slog.Int("request_id", requestID),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return nil
}

func (s *RequestService) CloseRequest(ctx context.Context, claims types.JWTClaims, requestID int) error {
	if _, err := s.checkUserCanEditRequest(ctx, claims, requestID); err != nil {
		return err
	}

	if err := s.requestRepo.CloseRequest(ctx, requestID); err != nil {
		s.logger.Error("error when trying to close request",
			slog.Int("request_id", requestID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("closed request successfully",
		slog.Int("request_id", requestID),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return nil
}

func (s *RequestService) CancelRequest(ctx context.Context, claims types.JWTClaims, requestID int) error {
	req, err := s.requestRepo.GetRequestByID(ctx, requestID)
	if err != nil {
		s.logger.Error("error when trying to retrieve request for cancelling it",
			slog.Int("request_id", requestID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return err
	}

	status := ComputeStatus(req.LastUploadedAt, req.NextDueAt, req.ExpectedDocuments)
	if status != types.StatusPending {
		s.logger.Warn("attempt to cancel non-pending request",
			slog.Int("request_id", requestID),
			slog.Int("jwt_user_id", claims.UserID),
		)
		return errors.ErrBadRequest{Msg: "A request cannot be closed if the status is not 'pending'."}
	}

	_, err = s.checkUserIsParticipantOfRequest(ctx, claims, requestID)
	if err != nil {
		s.logger.Error("user tried to cancel request that he does not own",
			slog.Int("request_id", requestID),
			slog.Int("jwt_user_id", claims.UserID),
		)
		return errors.ErrForbidden{Msg: "You are not allowed to cancel this request."}
	}

	err = s.requestRepo.CancelRequest(ctx, requestID)
	if err != nil {
		s.logger.Error("error when trying to cancel request",
			slog.Int("request_id", requestID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return err
	}

	return nil
}

func (service *RequestService) GetArchivedRequests(ctx context.Context, claims types.JWTClaims) ([]models.RequestDTORead, error) {
	if !claims.IsAdmin() && !claims.IsDepartmentMember() {
		service.logger.Warn("unauthorized attempt to get archived requests",
			slog.Int("jwt_user_id", claims.UserID),
		)
		return nil, errors.ErrForbidden{Msg: "Only admins and department members can view archived requests."}
	}

	reqs, err := service.fetchArchivedRequests(ctx, claims)
	if err != nil {
		service.logger.Error("failed to get archived requests",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	for i := range reqs {
		reqs[i].Status = ComputeStatus(reqs[i].LastUploadedAt, reqs[i].NextDueAt, reqs[i].ExpectedDocuments)
	}

	service.logger.Info("fetched archived requests successfully",
		slog.Int("jwt_user_id", claims.UserID),
	)
	return reqs, nil
}

func (service *RequestService) fetchArchivedRequests(ctx context.Context, claims types.JWTClaims) ([]models.RequestDTORead, error) {
	if claims.IsAdmin() {
		return service.requestRepo.GetArchivedRequests(ctx, nil)
	}
	return service.requestRepo.GetArchivedRequestsByDepartment(ctx, *claims.DepartmentID, nil)
}

func (service *RequestService) GetCancelledRequests(ctx context.Context, claims types.JWTClaims) ([]models.RequestDTORead, error) {
	if !claims.IsAdmin() && !claims.IsDepartmentMember() {
		service.logger.Warn("unauthorized attempt to get cancelled requests",
			slog.Int("jwt_user_id", claims.UserID),
		)
		return nil, errors.ErrForbidden{Msg: "Only admins and department members can view cancelled requests."}
	}

	reqs, err := service.fetchCancelledRequests(ctx, claims)
	if err != nil {
		service.logger.Error("failed to get cancelled requests",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	for i := range reqs {
		reqs[i].Status = ComputeStatus(reqs[i].LastUploadedAt, reqs[i].NextDueAt, reqs[i].ExpectedDocuments)
	}

	service.logger.Info("fetched cancelled requests successfully",
		slog.Int("jwt_user_id", claims.UserID),
	)
	return reqs, nil
}

func (service *RequestService) fetchCancelledRequests(ctx context.Context, claims types.JWTClaims) ([]models.RequestDTORead, error) {
	if claims.IsAdmin() {
		return service.requestRepo.GetCancelledRequests(ctx, nil)
	}
	return service.requestRepo.GetCancelledRequestsByDepartment(ctx, *claims.DepartmentID, nil)
}

func (s *RequestService) AddDocument(
	ctx context.Context,
	claims types.JWTClaims,
	requestID int,
	expectedDocID int,
	fileName string,
	fileSize int64,
	contentType string,
	content io.Reader,
) (*int, error) {
	if err := ValidateFileInfo(fileName, fileSize); err != nil {
		return nil, err
	}

	if _, err := s.checkUserIsParticipantOfRequest(ctx, claims, requestID); err != nil {
		return nil, err
	}

	s.logger.Info("attempting file upload",
		slog.Int("jwt_user_id", claims.UserID),
		slog.Int("request_id", requestID),
		slog.String("file_name", fileName),
	)

	cleanFileName := filepath.Base(fileName)
	s3Key := s.fileStorage.GenerateS3Key(fileName, fmt.Sprintf("/requests/%d", requestID))

	result, err := s.fileStorage.UploadFile(ctx, s3Key, content, contentType)
	if err != nil {
		s.logger.Error("s3 upload failed",
			slog.String("key", s3Key),
			slog.Any("error", err),
		)
		return nil, errors.ErrBadGateway{Msg: fmt.Sprintf("Failed to upload to S3. %v", err)}
	}

	fileModel := models.Document{
		RequestID:          requestID,
		ExpectedDocumentID: expectedDocID,
		FileName:           cleanFileName,
		FilePath:           s3Key,
		MimeType:           &contentType,
		FileSize:           &fileSize,
		S3VersionID:        result.VersionId,
		UploadedAt:         time.Now(),
		UploadedBy:         &claims.UserID,
	}

	id, err := s.requestRepo.AddDocument(ctx, fileModel)
	if err != nil {
		s.logger.Warn("metadata save failed, starting cleanup",
			slog.String("key", s3Key),
			slog.Any("error", err),
		)

		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cleanupErr := s.fileStorage.DeleteFile(cleanupCtx, fileModel.FilePath, fileModel.S3VersionID)
		if cleanupErr != nil {
			return nil, errors.ErrBadGateway{Msg: fmt.Sprintf("Failed to cleanup S3 object after DB failure. Key: %s, Version: %s, Error: %v\n", s3Key, *result.VersionId, cleanupErr)}
		}

		return nil, errors.ErrInternalServerError{Msg: fmt.Sprintf("Metadata save failed, file removed from storage. %v", err)}
	}

	// plain members uploading triggers status update
	if claims.Role == types.RoleMember && !claims.IsDepartmentMember() {
		s.requestRepo.SetFileUploaded(ctx, requestID)
		if expectedDocID != 0 {
			s.expectedDocRepo.UpdateExpectedDocumentStatus(ctx, expectedDocID, "uploaded", nil)
		}
	}

	s.logger.Info("file upload successful",
		slog.Int("file_id", id),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return &id, nil
}

func (s *RequestService) GetFilesByRequest(ctx context.Context, claims types.JWTClaims, requestID int) ([]models.DocumentDTORead, error) {
	if _, err := s.checkUserIsParticipantOfRequest(ctx, claims, requestID); err != nil {
		return nil, err
	}

	files, err := s.requestRepo.GetFilesByRequest(ctx, requestID)
	if err != nil {
		s.logger.Error("failed to fetch files",
			slog.Int("request_id", requestID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	s.logger.Info("files retrieved successfully",
		slog.Int("request_id", requestID),
		slog.Int("jwt_user_id", claims.UserID),
		slog.Int("count", len(files)),
	)
	return files, nil
}

func (s *RequestService) GetFilePresignedURL(ctx context.Context, claims types.JWTClaims, fileID int) (*string, error) {
	file, err := s.requestRepo.GetFileByID(ctx, fileID)
	if err != nil {
		s.logger.Error("could not fetch file by id",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("file_id", fileID),
			slog.Any("error", err),
		)
		return nil, err
	}

	if _, err := s.checkUserIsParticipantOfRequest(ctx, claims, file.RequestID); err != nil {
		return nil, err
	}

	presignedURL, err := s.fileStorage.GeneratePresignedURL(ctx, file.FilePath, file.S3VersionID, 15*time.Minute)
	if err != nil {
		s.logger.Error("s3 presign failed",
			slog.Int("file_id", fileID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}
	return &presignedURL, nil
}

func (s *RequestService) GetExamplePresignedURL(ctx context.Context, claims types.JWTClaims, expectedDocID int) (*string, error) {
	expectedDoc, err := s.expectedDocRepo.GetExpectedDocumentByID(ctx, expectedDocID)
	if err != nil {
		return nil, errors.ErrNotFound{Msg: "Expected document not found."}
	}

	if expectedDoc.ExampleFilePath == nil {
		return nil, errors.ErrNotFound{Msg: "This document has no example file."}
	}

	if _, err := s.checkUserIsParticipantOfRequest(ctx, claims, expectedDoc.RequestID); err != nil {
		return nil, err
	}

	presignedURL, err := s.fileStorage.GeneratePresignedURL(ctx, *expectedDoc.ExampleFilePath, nil, 15*time.Minute)
	if err != nil {
		s.logger.Error("s3 presign failed for example file",
			slog.Int("expected_doc_id", expectedDocID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}
	return &presignedURL, nil
}
