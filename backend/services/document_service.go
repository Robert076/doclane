package services

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types/errors"
)

type DocumentService struct {
	documentRepo    repositories.IDocumentRepository
	userRepo        repositories.IUserRepository
	expectedDocRepo repositories.IExpectedDocumentRepository
	txManager       repositories.ITxManager
	logger          *slog.Logger
	fileStorage     *FileStorageService
}

func NewDocumentService(documentRepo repositories.IDocumentRepository, userRepo repositories.IUserRepository, expectedDocRepo repositories.IExpectedDocumentRepository, txManager repositories.ITxManager, logger *slog.Logger, fileStorage *FileStorageService) *DocumentService {
	return &DocumentService{
		documentRepo:    documentRepo,
		userRepo:        userRepo,
		expectedDocRepo: expectedDocRepo,
		txManager:       txManager,
		logger:          logger,
		fileStorage:     fileStorage,
	}
}

func (service *DocumentService) AddDocumentRequest(
	ctx context.Context,
	jwtUserId int,
	dto models.DocumentRequestDTOCreate,
) (int, error) {
	if err := ValidateRequestInput(dto); err != nil {
		service.logger.Warn("document request create failed because it did not pass validations",
			slog.Int("user_id", jwtUserId),
			slog.Any("error", err))
		return 0, err
	}

	client, err := service.userRepo.GetUserByID(ctx, dto.ClientID)
	if err != nil {
		service.logger.Warn("client lookup failed for document request",
			slog.Int("client_id", dto.ClientID),
			slog.Int("requested_by", jwtUserId),
			slog.Any("error", err),
		)
		return 0, errors.ErrNotFound{Msg: "Client not found."}
	}

	if client.ProfessionalID == nil || *client.ProfessionalID != jwtUserId {
		service.logger.Warn("unauthorized attempt to add request to unassigned client",
			slog.Int("professional_id", jwtUserId),
			slog.Int("client_id", dto.ClientID),
		)
		return 0, errors.ErrForbidden{Msg: "This client is not assigned to you."}
	}

	nextDueAt := ComputeNextDueAt(dto.DueDate, dto.RecurrenceCron)

	req := models.DocumentRequest{
		ProfessionalID: jwtUserId,
		DocumentRequestBase: models.DocumentRequestBase{
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

	var id int
	err = service.txManager.WithTx(ctx, func(tx *sql.Tx) error {
		var txErr error
		id, txErr = service.documentRepo.AddDocumentRequestWithTx(ctx, req, tx)
		if txErr != nil {
			return txErr
		}

		for _, ed := range dto.ExpectedDocuments {
			ed.DocumentRequestID = id
			ed.IsUploaded = false
			if txErr = service.expectedDocRepo.AddExpectedDocumentToRequestWithTx(ctx, tx, ed); txErr != nil {
				return txErr
			}
		}
		return nil
	})
	if err != nil {
		service.logger.Error("failed to create document request",
			slog.Any("error", err),
			slog.Int("professional_id", jwtUserId),
			slog.Int("client_id", dto.ClientID),
		)
		return 0, err
	}

	service.logger.Info("document request created",
		slog.Int("request_id", id),
		slog.Int("expected_documents", len(dto.ExpectedDocuments)),
	)
	return id, nil
}

func (service *DocumentService) GetDocumentRequestByID(
	ctx context.Context,
	jwtUserId int,
	id int,
) (models.DocumentRequestDTORead, error) {
	req, err := service.documentRepo.GetDocumentRequestByID(ctx, id)
	if err != nil {
		service.logger.Error("failed to get document request by id",
			slog.Int("request_id", id),
			slog.Any("error", err),
		)
		return models.DocumentRequestDTORead{}, err
	}

	if req.ProfessionalID != jwtUserId && req.ClientID != jwtUserId {
		service.logger.Warn("unauthorized access attempt to document request",
			slog.Int("user_id", jwtUserId),
			slog.Int("request_id", id),
		)
		return models.DocumentRequestDTORead{}, errors.ErrForbidden{Msg: fmt.Sprintf("User with id %v is not allowed to access document request with id %v", jwtUserId, req.ID)}
	}

	expectedDocs, err := service.expectedDocRepo.GetExpectedDocumentsByRequest(ctx, id)
	if err != nil {
		service.logger.Error("failed to get expected documents for request",
			slog.Int("request_id", id),
			slog.Any("error", err),
		)
		return models.DocumentRequestDTORead{}, err
	}
	req.ExpectedDocuments = expectedDocs

	req.Status = ComputeStatus(req.LastUploadedAt, req.NextDueAt)

	return req, nil
}

func (service *DocumentService) GetDocumentRequestsByProfessional(
	ctx context.Context,
	jwtUserId int,
	search *string,
) ([]models.DocumentRequestDTORead, error) {
	return service.getDocumentRequestsByRole(ctx, jwtUserId, "PROFESSIONAL", search, service.documentRepo.GetDocumentRequestsByProfessional)
}

func (service *DocumentService) GetDocumentRequestsByClient(
	ctx context.Context,
	jwtUserId int,
	search *string,
) ([]models.DocumentRequestDTORead, error) {
	return service.getDocumentRequestsByRole(ctx, jwtUserId, "CLIENT", search, service.documentRepo.GetDocumentRequestsByClient)
}

func (service *DocumentService) getDocumentRequestsByRole(
	ctx context.Context,
	jwtUserId int,
	requiredRole string,
	search *string,
	fetchFunc func(context.Context, int, *string) ([]models.DocumentRequestDTORead, error),
) ([]models.DocumentRequestDTORead, error) {
	user, err := service.userRepo.GetUserByID(ctx, jwtUserId)
	if err != nil {
		service.logger.Error("failed to fetch user for document requests",
			slog.Int("user_id", jwtUserId),
			slog.Any("error", err),
		)
		return nil, err
	}

	if user.Role != requiredRole {
		service.logger.Warn("user tried to access other role's endpoint for document requests",
			slog.Int("user_id", jwtUserId),
			slog.String("role", user.Role))
		return nil, errors.ErrForbidden{Msg: fmt.Sprintf("This is a %s endpoint.", requiredRole)}
	}

	reqs, err := fetchFunc(ctx, jwtUserId, search)
	if err != nil {
		service.logger.Error("failed to fetch document requests from repo",
			slog.Int("client_id", jwtUserId),
			slog.Any("error", err),
		)
		return nil, err
	}

	for i := range reqs {
		reqs[i].Status = ComputeStatus(reqs[i].LastUploadedAt, reqs[i].NextDueAt)
	}

	return reqs, nil
}

func (service *DocumentService) PatchDocumentRequest(
	ctx context.Context,
	jwtUserID int,
	requestID int,
	updatedDTO models.DocumentRequestDTOPatch,
) error {
	if err := ValidatePatchDTO(updatedDTO); err != nil {
		service.logger.Warn("patch validation failed",
			slog.Int("user_id", jwtUserID),
			slog.Int("request_id", requestID),
			slog.Any("error", err),
		)
		return err
	}

	req, err := service.documentRepo.GetDocumentRequestByID(ctx, requestID)
	if err != nil {
		service.logger.Error("failed to get document request for patch",
			slog.Int("request_id", requestID),
			slog.Any("error", err),
		)
		return err
	}

	if req.ProfessionalID != jwtUserID {
		service.logger.Warn("unauthorized patch attempt",
			slog.Int("user_id", jwtUserID),
			slog.Int("request_id", requestID),
			slog.Int("actual_professional_id", req.ProfessionalID),
		)
		return errors.ErrForbidden{Msg: "Forbidden."}
	}

	if err := service.documentRepo.UpdateDocumentRequestTitle(ctx, requestID, updatedDTO.Title); err != nil {
		service.logger.Error("failed to update document request title",
			slog.Int("request_id", requestID),
			slog.String("new_title", updatedDTO.Title),
			slog.Any("error", err),
		)
		return err
	}

	service.logger.Info("document request patched successfully",
		slog.Int("request_id", requestID),
		slog.String("new_title", updatedDTO.Title),
	)

	return nil
}

func (service *DocumentService) CloseRequest(
	ctx context.Context,
	jwtUserID int,
	requestID int,
) error {
	req, err := service.documentRepo.GetDocumentRequestByID(ctx, requestID)
	if err != nil {
		return err
	}

	if req.ProfessionalID != jwtUserID {
		return errors.ErrForbidden{Msg: "You are not allowed to close this request."}
	}

	if err := service.documentRepo.CloseDocumentRequest(ctx, requestID); err != nil {
		return err
	}

	return nil
}

func (service *DocumentService) AddDocumentFile(
	ctx context.Context,
	jwtUserId int,
	requestID int,
	expectedDocID int,
	fileName string,
	fileSize int64,
	contentType string,
	content io.Reader,
) (int, error) {
	if err := ValidateFileInfo(fileName, fileSize); err != nil {
		return 0, err
	}

	service.logger.Info("attempting file upload",
		slog.Int("user_id", jwtUserId),
		slog.Int("request_id", requestID),
		slog.String("file_name", fileName),
	)

	docReq, err := service.documentRepo.GetDocumentRequestByID(ctx, requestID)
	if err != nil {
		return 0, errors.ErrNotFound{Msg: fmt.Sprintf("Document request not found. %v", err)}
	}

	if docReq.ClientID != jwtUserId && docReq.ProfessionalID != jwtUserId {
		return 0, errors.ErrForbidden{Msg: fmt.Sprintf("User with id %v is not allowed to modify document request with id %v.", jwtUserId, requestID)}
	}

	cleanFileName := filepath.Base(fileName)
	s3Key := generateS3Key(fileName, requestID)

	result, err := service.fileStorage.UploadFile(ctx, s3Key, content, contentType)
	if err != nil {
		service.logger.Error("s3 upload failed",
			slog.String("key", s3Key),
			slog.Any("error", err),
		)
		return 0, errors.ErrBadGateway{Msg: fmt.Sprintf("Failed to upload to S3. %v", err)}
	}

	fileModel := models.DocumentFile{
		DocumentRequestID:  requestID,
		ExpectedDocumentID: expectedDocID,
		FileName:           cleanFileName,
		FilePath:           s3Key,
		MimeType:           &contentType,
		FileSize:           &fileSize,
		S3VersionID:        result.VersionId,
		UploadedAt:         time.Now(),
		UploadedBy:         &jwtUserId,
	}

	id, err := service.documentRepo.AddDocumentFile(ctx, fileModel)
	if err != nil {
		service.logger.Warn("metadata save failed, starting cleanup",
			slog.String("key", s3Key),
			slog.Any("error", err),
		)

		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cleanupErr := service.fileStorage.DeleteFile(cleanupCtx, fileModel.FilePath, fileModel.S3VersionID)
		if cleanupErr != nil {
			return 0, errors.ErrBadGateway{Msg: fmt.Sprintf("Failed to cleanup S3 object after DB failure. Key: %s, Version: %s, Error: %v\n", s3Key, *result.VersionId, cleanupErr)}
		}

		return 0, errors.ErrInternalServerError{Msg: fmt.Sprintf("Metadata save failed, file removed from storage. %v", err)}
	}

	uploadedFile, err := service.documentRepo.GetFileByIDExtended(ctx, id)
	if err != nil {
		service.logger.Error("error getting uploaded file",
			slog.Int("id", uploadedFile.ID),
			slog.Any("err", err),
		)
		return 0, errors.ErrInternalServerError{Msg: fmt.Sprintf("Error getting uploaded file: %v", err)}
	}

	if uploadedFile.AuthorRole == "CLIENT" {
		// Maybe the professional uploads an example, the request should not be marked as finished
		service.documentRepo.SetFileUploaded(ctx, requestID)
	}

	service.logger.Info("file upload successful", slog.Int("file_id", id))
	return id, nil
}

func (service *DocumentService) GetFilesByRequest(
	ctx context.Context,
	jwtUserId int,
	requestID int,
) ([]models.DocumentFileDTORead, error) {
	docReq, err := service.documentRepo.GetDocumentRequestByID(ctx, requestID)
	if err != nil {
		service.logger.Error("failed to find document request", slog.Int("request_id", requestID), slog.Any("error", err))
		return nil, err
	}

	if docReq.ProfessionalID != jwtUserId && docReq.ClientID != jwtUserId {
		service.logger.Warn("unauthorized access attempt", slog.Int("user_id", jwtUserId), slog.Int("request_id", requestID))
		return nil, errors.ErrForbidden{Msg: "You are not authorized to view files for this request."}
	}

	files, err := service.documentRepo.GetFilesByRequest(ctx, requestID)
	if err != nil {
		service.logger.Error("failed to fetch files", slog.Int("request_id", requestID), slog.Any("error", err))
		return nil, err
	}

	service.logger.Info("files retrieved successfully", slog.Int("request_id", requestID), slog.Int("count", len(files)))
	return files, nil
}

func (service *DocumentService) GetFilePresignedURL(
	ctx context.Context,
	jwtUserId int,
	fileID int,
) (string, error) {
	file, err := service.documentRepo.GetFileByID(ctx, fileID)
	if err != nil {
		return "", err
	}

	docReq, err := service.documentRepo.GetDocumentRequestByID(ctx, file.DocumentRequestID)
	if err != nil {
		return "", err
	}

	if docReq.ProfessionalID != jwtUserId && docReq.ClientID != jwtUserId {
		return "", errors.ErrForbidden{Msg: "No access to this file."}
	}

	presignedURL, err := service.fileStorage.GeneratePresignedURL(ctx, file.FilePath, file.S3VersionID, 15*time.Minute)
	if err != nil {
		service.logger.Error("s3 presign failed",
			slog.Int("file_id", fileID),
			slog.Any("error", err))
		return "", err
	}
	return presignedURL, nil
}
