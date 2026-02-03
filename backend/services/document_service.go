package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type DocumentService struct {
	documentRepo repositories.IDocumentRepository
	userRepo     repositories.IUserRepository
	s3Client     *s3.Client
	bucket       string
	logger       *slog.Logger
}

func NewDocumentService(documentRepo repositories.IDocumentRepository, userRepo repositories.IUserRepository, s3Client *s3.Client, bucket string, logger *slog.Logger) *DocumentService {
	return &DocumentService{
		documentRepo: documentRepo,
		userRepo:     userRepo,
		s3Client:     s3Client,
		bucket:       bucket,
		logger:       logger,
	}
}

func (service *DocumentService) AddDocumentRequest(
	ctx context.Context,
	jwtUserId int,
	dto models.DocumentRequestDTOCreate,
) (int, error) {
	if err := service.validateRequestInput(dto.Title, dto.DueDate); err != nil {
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

	jwtUserIdStr := strconv.Itoa(jwtUserId)
	if client.ProfessionalID == nil || *client.ProfessionalID != jwtUserIdStr {
		service.logger.Warn("unauthorized attempt to add request to unassigned client",
			slog.Int("professional_id", jwtUserId),
			slog.Int("client_id", dto.ClientID),
		)
		return 0, errors.ErrForbidden{Msg: "This client is not assigned to you."}
	}

	nextDueAt := computeNextDueAt(dto.DueDate, dto.RecurrenceCron)
	req := models.DocumentRequest{
		ProfessionalID: jwtUserId,
		DocumentRequestBase: models.DocumentRequestBase{
			ClientID:       dto.ClientID,
			Title:          dto.Title,
			Description:    dto.Description,
			IsRecurring:    dto.IsRecurring,
			RecurrenceCron: dto.RecurrenceCron,
			NextDueAt:      nextDueAt,
			LastUploadedAt: dto.LastUploadedAt,
			DueDate:        dto.DueDate,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := service.documentRepo.AddDocumentRequest(ctx, req)
	if err != nil {
		service.logger.Error("failed to create document request",
			slog.Any("error", err),
			slog.Int("professional_id", jwtUserId),
			slog.Int("client_id", dto.ClientID),
		)
		return 0, err
	}

	service.logger.Info("document request created", slog.Int("request_id", id))
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

	req.Status = computeStatus(req.LastUploadedAt, req.NextDueAt)

	return req, nil
}

func (service *DocumentService) GetDocumentRequestsByProfessional(
	ctx context.Context,
	jwtUserId int,
) ([]models.DocumentRequestDTORead, error) {
	user, err := service.userRepo.GetUserByID(ctx, jwtUserId)
	if err != nil {
		service.logger.Error("failed to fetch professional for document requests",
			slog.Int("user_id", jwtUserId),
			slog.Any("error", err))
		return nil, err
	}

	if user.Role != "PROFESSIONAL" {
		service.logger.Warn("non-professional tried to access professional endpoint for document requests",
			slog.Int("user_id", jwtUserId),
			slog.String("role", user.Role))
		return nil, errors.ErrForbidden{Msg: "This is a professional endpoint."}
	}

	reqs, err := service.documentRepo.GetDocumentRequestsByProfessional(ctx, jwtUserId)
	if err != nil {
		service.logger.Error("failed to fetch professional document requests",
			slog.Int("user_id", jwtUserId),
			slog.Any("error", err),
		)
		return nil, err
	}

	for i := range reqs {
		reqs[i].Status = computeStatus(reqs[i].LastUploadedAt, reqs[i].NextDueAt)
	}

	return reqs, nil
}

func (service *DocumentService) GetDocumentRequestsByClient(
	ctx context.Context,
	jwtUserId int,
) ([]models.DocumentRequestDTORead, error) {
	user, err := service.userRepo.GetUserByID(ctx, jwtUserId)
	if err != nil {
		service.logger.Error("failed to fetch client for document requests",
			slog.Int("user_id", jwtUserId),
			slog.Any("error", err),
		)
		return nil, err
	}

	if user.Role != "CLIENT" {
		service.logger.Warn("non-client tried to access client endpoint for document requests",
			slog.Int("user_id", jwtUserId),
			slog.String("role", user.Role))
		return nil, errors.ErrForbidden{Msg: "This is a client endpoint."}
	}

	reqs, err := service.documentRepo.GetDocumentRequestsByClient(ctx, jwtUserId)
	if err != nil {
		service.logger.Error("failed to fetch document requests from repo",
			slog.Int("client_id", jwtUserId),
			slog.Any("error", err),
		)
		return nil, err
	}

	for i := range reqs {
		reqs[i].Status = computeStatus(reqs[i].LastUploadedAt, reqs[i].NextDueAt)
	}

	return reqs, nil
}

func (service *DocumentService) AddDocumentFile(
	ctx context.Context,
	jwtUserId int,
	requestID int,
	fileName string,
	fileSize int64,
	contentType string,
	content io.Reader,
) (int, error) {
	if err := service.validateFileInfo(fileName, fileSize); err != nil {
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
	uniqueID := uuid.New().String()

	s3Key := fmt.Sprintf("requests/%d/%s-%s", requestID, uniqueID, cleanFileName)

	result, err := service.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(service.bucket),
		Key:         aws.String(s3Key),
		Body:        content,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		service.logger.Error("s3 upload failed",
			slog.String("key", s3Key),
			slog.Any("error", err),
		)
		return 0, errors.ErrBadGateway{Msg: fmt.Sprintf("Failed to upload to S3. %v", err)}
	}

	fileModel := models.DocumentFile{
		DocumentRequestID: requestID,
		FileName:          cleanFileName,
		FilePath:          s3Key,
		MimeType:          &contentType,
		FileSize:          &fileSize,
		S3VersionID:       result.VersionId,
		UploadedAt:        time.Now(),
	}

	id, err := service.documentRepo.AddDocumentFile(ctx, fileModel)
	if err != nil {
		service.logger.Warn("metadata save failed, starting cleanup",
			slog.String("key", s3Key),
			slog.Any("error", err),
		)

		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, cleanupErr := service.s3Client.DeleteObject(cleanupCtx, &s3.DeleteObjectInput{
			Bucket:    aws.String(service.bucket),
			Key:       aws.String(s3Key),
			VersionId: result.VersionId,
		})

		if cleanupErr != nil {
			service.logger.Error("S3 CLEANUP FAILED - ZOMBIE FILE ALERT",
				slog.String("key", s3Key),
				slog.Any("error", cleanupErr),
			)

			return 0, errors.ErrBadGateway{Msg: fmt.Sprintf("Failed to cleanup S3 object after DB failure. Key: %s, Version: %s, Error: %v\n", s3Key, *result.VersionId, cleanupErr)}
		}

		return 0, errors.ErrInternalServerError{Msg: fmt.Sprintf("Metadata save failed, file removed from storage. %v", err)}
	}

	service.documentRepo.SetFileUploaded(ctx, requestID)

	service.logger.Info("file upload successful", slog.Int("file_id", id))
	return id, nil
}

func (service *DocumentService) GetFilesByRequest(
	ctx context.Context,
	jwtUserId int,
	requestID int,
) ([]models.DocumentFile, error) {
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

	presignClient := s3.NewPresignClient(service.s3Client)
	presignParams := &s3.GetObjectInput{
		Bucket:    aws.String(service.bucket),
		Key:       aws.String(file.FilePath),
		VersionId: file.S3VersionID,
	}

	presignedReq, err := presignClient.PresignGetObject(ctx, presignParams, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})

	if err != nil {
		service.logger.Error("s3 presign failed",
			slog.Int("file_id", fileID),
			slog.Any("error", err))
		return "", err
	}

	return presignedReq.URL, nil
}

func (service *DocumentService) validateRequestInput(title string, dueDate *time.Time) error {
	if len(title) < 3 || len(title) > 40 {
		return errors.ErrBadRequest{Msg: "Title must be between 3 and 40 characters."}
	}

	if dueDate != nil && dueDate.Before(time.Now()) {
		return errors.ErrBadRequest{Msg: "Due date cannot be in the past."}
	}

	return nil
}

func (service *DocumentService) validateFileInfo(fileName string, fileSize int64) error {
	if fileSize <= 0 {
		return errors.ErrBadRequest{Msg: "File is empty."}
	}

	const maxFileSize = 20 * 1024 * 1024
	if fileSize > maxFileSize {
		return errors.ErrBadRequest{Msg: "File size must be less than 20MB."}
	}

	allowedExtensions := map[string]bool{
		".pdf": true, ".jpg": true, ".jpeg": true, ".png": true, ".doc": true, ".docx": true,
	}
	ext := filepath.Ext(fileName)
	if !allowedExtensions[ext] {
		return errors.ErrBadRequest{Msg: fmt.Sprintf("Extension %s is not allowed.", ext)}
	}

	return nil
}

func computeNextDueAt(dueDate *time.Time, cronExpr *string) *time.Time {
	now := time.Now()

	if dueDate != nil {
		return dueDate
	}

	if cronExpr == nil || *cronExpr == "" {
		return nil
	}

	schedule, err := cron.ParseStandard(*cronExpr)
	if err != nil {
		return nil
	}

	next := schedule.Next(now)

	return &next
}

func computeStatus(lastUploadedAt *time.Time, nextDueAt *time.Time) string {
	now := time.Now()

	if nextDueAt == nil {
		return "pending"
	}

	if lastUploadedAt != nil && !lastUploadedAt.Before(*nextDueAt) {
		return "uploaded"
	}

	if now.After(*nextDueAt) {
		return "overdue"
	}

	return "pending"
}
