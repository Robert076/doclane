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
	"github.com/Robert076/doclane/backend/types/errors"
)

type RequestService struct {
	requestRepo     repositories.IRequestRepo
	userRepo        repositories.IUserRepo
	expectedDocRepo repositories.IExpectedDocumentRepo
	txManager       repositories.ITxManager
	fileStorage     IFileStorageService
	logger          *slog.Logger
}

func NewRequestService(requestRepo repositories.IRequestRepo, userRepo repositories.IUserRepo, expectedDocRepo repositories.IExpectedDocumentRepo, txManager repositories.ITxManager, logger *slog.Logger, fileStorage IFileStorageService) *RequestService {
	return &RequestService{
		requestRepo:     requestRepo,
		userRepo:        userRepo,
		expectedDocRepo: expectedDocRepo,
		txManager:       txManager,
		logger:          logger,
		fileStorage:     fileStorage,
	}
}

func (service *RequestService) AddRequest(
	ctx context.Context,
	jwtUserID int,
	dto models.RequestDTOCreate,
) (*int, error) {
	if err := ValidateRequestInput(dto); err != nil {
		service.logger.Warn("document request create failed because it did not pass validations",
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err))
		return nil, err
	}

	client, err := service.userRepo.GetUserByID(ctx, dto.ClientID)
	if err != nil {
		service.logger.Warn("client lookup failed for document request",
			slog.Int("client_id", dto.ClientID),
			slog.Int("requested_by", jwtUserID),
			slog.Any("error", err),
		)
		return nil, errors.ErrNotFound{Msg: "Client not found."}
	}

	if client.ProfessionalID == nil || *client.ProfessionalID != jwtUserID {
		service.logger.Warn("unauthorized attempt to add request to unassigned client",
			slog.Int("professional_id", jwtUserID),
			slog.Int("client_id", dto.ClientID),
		)
		return nil, errors.ErrForbidden{Msg: "This client is not assigned to you."}
	}

	nextDueAt := ComputeNextDueAt(dto.DueDate, dto.RecurrenceCron)

	req := models.Request{
		ProfessionalID: jwtUserID,
		RequestBase: models.RequestBase{
			ClientID:       dto.ClientID,
			Title:          dto.Title,
			Description:    dto.Description,
			IsRecurring:    dto.IsRecurring,
			RecurrenceCron: dto.RecurrenceCron,
			IsScheduled:    dto.IsScheduled,
			ScheduledFor:   dto.ScheduledFor,
			NextDueAt:      nextDueAt,
			LastUploadedAt: dto.LastUploadedAt,
			DueDate:        dto.DueDate,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	uploadedExamples, err := service.getUploadedExamples(ctx, dto)
	if err != nil {
		return nil, err
	}

	expectedDocs := getExpectedDocuments(dto, uploadedExamples)

	id, err := service.createRequestTransaction(ctx, req, expectedDocs)
	if err != nil {
		service.logger.Error("failed to create document request",
			slog.Any("error", err),
			slog.Int("professional_id", jwtUserID),
			slog.Int("client_id", dto.ClientID),
		)

		service.removeUploadedExamples(ctx, uploadedExamples)
		return nil, err
	}

	service.logger.Info("document request created",
		slog.Int("request_id", *id),
		slog.Int("expected_documents", len(dto.ExpectedDocuments)),
	)
	return id, nil
}

func (s *RequestService) GetRequestByID(
	ctx context.Context,
	jwtUserID int,
	requestID int,
) (*models.RequestDTORead, error) {
	req, err := s.checkUserIsParticipantOfRequest(ctx, jwtUserID, requestID)
	if err != nil {
		return nil, err
	}

	expectedDocs, err := s.expectedDocRepo.GetExpectedDocumentsByRequest(ctx, requestID)
	if err != nil {
		s.logger.Error("failed to get expected documents for request",
			slog.Int("request_id", requestID),
			slog.Any("error", err),
		)
		return nil, err
	}
	req.ExpectedDocuments = expectedDocs

	req.Status = ComputeStatus(req.LastUploadedAt, req.NextDueAt, req.ExpectedDocuments)

	return req, nil
}

func (service *RequestService) GetRequestsByProfessional(
	ctx context.Context,
	jwtUserID int,
	search *string,
) ([]models.RequestDTORead, error) {
	return service.getRequestsByRole(ctx, jwtUserID, "PROFESSIONAL", search, service.requestRepo.GetRequestsByProfessionalWithExpectedDocs)
}

func (service *RequestService) GetRequestsByClient(
	ctx context.Context,
	jwtUserID int,
	search *string,
) ([]models.RequestDTORead, error) {
	return service.getRequestsByRole(ctx, jwtUserID, "CLIENT", search, service.requestRepo.GetRequestsByClientWithExpectedDocs)
}

func (service *RequestService) getRequestsByRole(
	ctx context.Context,
	jwtUserID int,
	requiredRole string,
	search *string,
	fetchFunc func(context.Context, int, *string) ([]models.RequestDTORead, error),
) ([]models.RequestDTORead, error) {
	user, err := service.userRepo.GetUserByID(ctx, jwtUserID)
	if err != nil {
		service.logger.Error("failed to fetch user for document requests",
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	if user.Role != requiredRole {
		service.logger.Warn("user tried to access other role's endpoint for document requests",
			slog.Int("user_id", jwtUserID),
			slog.String("role", user.Role))
		return nil, errors.ErrForbidden{Msg: fmt.Sprintf("This is a %s endpoint.", requiredRole)}
	}

	reqs, err := fetchFunc(ctx, jwtUserID, search)
	if err != nil {
		service.logger.Error("failed to fetch document requests from repo",
			slog.Int("client_id", jwtUserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	for i := range reqs {
		reqs[i].Status = ComputeStatus(reqs[i].LastUploadedAt, reqs[i].NextDueAt, reqs[i].ExpectedDocuments)
	}

	return reqs, nil
}

func (s *RequestService) PatchRequest(
	ctx context.Context,
	jwtUserID int,
	requestID int,
	updatedDTO models.RequestDTOPatch,
) error {
	if err := ValidatePatchDTO(updatedDTO); err != nil {
		s.logger.Warn("patch validation failed",
			slog.Int("user_id", jwtUserID),
			slog.Int("request_id", requestID),
			slog.Any("error", err),
		)
		return err
	}

	if _, err := s.checkUserIsProfessionalOfRequest(ctx, jwtUserID, requestID); err != nil {
		return err
	}

	if err := s.requestRepo.UpdateRequestTitle(ctx, requestID, updatedDTO.Title); err != nil {
		s.logger.Error("failed to update document request title",
			slog.Int("request_id", requestID),
			slog.String("new_title", updatedDTO.Title),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("document request patched successfully",
		slog.Int("request_id", requestID),
		slog.String("new_title", updatedDTO.Title),
	)

	return nil
}

func (s *RequestService) ReopenRequest(
	ctx context.Context,
	jwtUserID int,
	requestID int,
) error {
	if _, err := s.checkUserIsProfessionalOfRequest(ctx, jwtUserID, requestID); err != nil {
		return err
	}

	return s.requestRepo.ReopenRequest(ctx, requestID)
}

func (s *RequestService) CloseRequest(
	ctx context.Context,
	jwtUserID int,
	requestID int,
) error {
	if _, err := s.checkUserIsProfessionalOfRequest(ctx, jwtUserID, requestID); err != nil {
		return err
	}

	return s.requestRepo.CloseRequest(ctx, requestID)
}

func (s *RequestService) AddDocument(
	ctx context.Context,
	jwtUserID int,
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

	if _, err := s.checkUserIsParticipantOfRequest(ctx, jwtUserID, requestID); err != nil {
		return nil, err
	}
	s.logger.Info("attempting file upload",
		slog.Int("user_id", jwtUserID),
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
		UploadedBy:         &jwtUserID,
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

	uploadedFile, err := s.requestRepo.GetFileByIDExtended(ctx, id)
	if err != nil {
		s.logger.Error("error getting uploaded file",
			slog.Int("id", id),
			slog.Any("err", err),
		)
		return nil, errors.ErrInternalServerError{Msg: fmt.Sprintf("Error getting uploaded file: %v", err)}
	}

	if uploadedFile.AuthorRole == "CLIENT" {
		s.requestRepo.SetFileUploaded(ctx, requestID)
		if expectedDocID != 0 {
			s.expectedDocRepo.UpdateExpectedDocumentStatus(ctx, expectedDocID, "uploaded", nil)
		}
	}

	s.logger.Info("file upload successful", slog.Int("file_id", id))
	return &id, nil
}

func (s *RequestService) GetFilesByRequest(
	ctx context.Context,
	jwtUserID int,
	requestID int,
) ([]models.DocumentDTORead, error) {
	if _, err := s.checkUserIsParticipantOfRequest(ctx, jwtUserID, requestID); err != nil {
		return nil, err
	}

	files, err := s.requestRepo.GetFilesByRequest(ctx, requestID)
	if err != nil {
		s.logger.Error("failed to fetch files", slog.Int("request_id", requestID), slog.Any("error", err))
		return nil, err
	}

	s.logger.Info("files retrieved successfully", slog.Int("request_id", requestID), slog.Int("count", len(files)))
	return files, nil
}

func (s *RequestService) GetFilePresignedURL(
	ctx context.Context,
	jwtUserID int,
	fileID int,
) (*string, error) {
	file, err := s.requestRepo.GetFileByID(ctx, fileID)
	if err != nil {
		s.logger.Error("could not fetch file by id",
			slog.Int("user_id", jwtUserID),
			slog.Int("file_id", fileID),
			slog.Any("error", err),
		)
		return nil, err
	}

	if _, err := s.checkUserIsParticipantOfRequest(ctx, jwtUserID, file.RequestID); err != nil {
		return nil, err
	}

	presignedURL, err := s.fileStorage.GeneratePresignedURL(ctx, file.FilePath, file.S3VersionID, 15*time.Minute)
	if err != nil {
		s.logger.Error("s3 presign failed",
			slog.Int("file_id", fileID),
			slog.Any("error", err))
		return nil, err
	}
	return &presignedURL, nil
}

func (s *RequestService) GetExamplePresignedURL(
	ctx context.Context,
	jwtUserID int,
	expectedDocID int,
) (*string, error) {
	expectedDoc, err := s.expectedDocRepo.GetExpectedDocumentByID(ctx, expectedDocID)
	if err != nil {
		return nil, errors.ErrNotFound{Msg: "Expected document not found."}
	}

	if expectedDoc.ExampleFilePath == nil {
		return nil, errors.ErrNotFound{Msg: "This document has no example file."}
	}

	if _, err := s.checkUserIsParticipantOfRequest(ctx, jwtUserID, expectedDoc.RequestID); err != nil {
		return nil, err
	}

	presignedURL, err := s.fileStorage.GeneratePresignedURL(ctx, *expectedDoc.ExampleFilePath, nil, 15*time.Minute)
	if err != nil {
		s.logger.Error("s3 presign failed for example file",
			slog.Int("expected_doc_id", expectedDocID),
			slog.Any("error", err),
		)
		return nil, err
	}
	return &presignedURL, nil
}
