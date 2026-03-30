package services

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/robfig/cron/v3"
)

func ValidateRequestInput(dto models.RequestDTOCreate) error {
	if len(dto.Title) < 3 || len(dto.Title) > 30 {
		return errors.ErrBadRequest{Msg: "Title must be between 3 and 40 characters."}
	}

	if dto.DueDate != nil && dto.DueDate.Before(time.Now()) {
		return errors.ErrBadRequest{Msg: "Due date cannot be in the past."}
	}

	if dto.IsRecurring == true && dto.RecurrenceCron == nil {
		return errors.ErrUnprocessableContent{Msg: "A request that is marked as recurring (is_recurring = true) should have a recurrence_cron field that is not null."}
	}

	if dto.IsScheduled == true && dto.ScheduledFor == nil {
		return errors.ErrUnprocessableContent{Msg: "A request that is marked as scheduled (is_scheduled = true) should have a scheduled_for field that is not null."}
	}

	if len(dto.ExpectedDocuments) == 0 {
		return errors.ErrUnprocessableContent{Msg: "Must provide at least 1 expected document for a new request."}
	}

	return nil
}

func ComputeStatus(lastUploadedAt *time.Time, nextDueAt *time.Time, expectedDocs []models.ExpectedDocument) string {
	now := time.Now()

	allUploaded := len(expectedDocs) > 0 && func() bool {
		for _, ed := range expectedDocs {
			if ed.Status != "uploaded" && ed.Status != "accepted" {
				return false
			}
		}
		return true
	}()

	if nextDueAt == nil {
		if allUploaded {
			return "uploaded"
		}
		if lastUploadedAt != nil {
			return "uploaded"
		}
		return "pending"
	}

	if now.After(*nextDueAt) {
		if allUploaded {
			return "uploaded"
		}
		return "overdue"
	}

	if allUploaded {
		return "uploaded"
	}

	return "pending"
}

func ValidatePatchDTO(dto models.RequestDTOPatch) error {
	if len(dto.Title) < 3 || len(dto.Title) > 30 {
		return errors.ErrBadRequest{Msg: "New title is too short or too long. Minimum 3 characters, maximum 30 characters."}
	}

	return nil
}

func ValidateRequestTemplateInput(template models.RequestTemplate) error {
	fmt.Print(template.Title)
	if len(template.Title) < 3 || len(template.Title) > 30 {
		return errors.ErrBadRequest{Msg: "Title must be between 3 and 30 characters."}
	}

	if template.IsRecurring && (template.RecurrenceCron == nil || *template.RecurrenceCron == "") {
		return errors.ErrUnprocessableContent{Msg: "A template marked as recurring must have a recurrence_cron field."}
	}

	if template.RecurrenceCron != nil && *template.RecurrenceCron != "" {
		if _, err := cron.ParseStandard(*template.RecurrenceCron); err != nil {
			return errors.ErrBadRequest{Msg: "Invalid recurrence_cron format."}
		}
	}

	return nil
}

func ComputeNextDueAt(dueDate *time.Time, cronExpr *string) *time.Time {
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

func (s *RequestService) checkUserIsProfessionalOfRequest(ctx context.Context, jwtUserID int, requestID int) (*models.RequestDTORead, error) {
	req, err := s.requestRepo.GetRequestByID(ctx, requestID)
	if err != nil {
		s.logger.Error("error getting request from db",
			slog.Int("user_id", jwtUserID),
			slog.Int("request_id", requestID),
			slog.Any("error", err),
		)
		return nil, err
	}

	if req.ProfessionalID != jwtUserID {
		return nil, errors.ErrForbidden{Msg: "You don't have access to the request."}
	}

	return &req, nil
}

func (s *RequestService) checkUserIsParticipantOfRequest(ctx context.Context, jwtUserID int, requestID int) (*models.RequestDTORead, error) {
	req, err := s.requestRepo.GetRequestByID(ctx, requestID)
	if err != nil {
		s.logger.Error("error getting request from db",
			slog.Int("user_id", jwtUserID),
			slog.Int("request_id", requestID),
			slog.Any("error", err),
		)
		return nil, err
	}

	if req.ProfessionalID != jwtUserID && req.ClientID != jwtUserID {
		s.logger.Warn("unauthorized access attempted for request",
			slog.Int("user_id", jwtUserID),
			slog.Int("request_id", requestID),
		)
		return nil, errors.ErrForbidden{Msg: "You don't have access to this request."}
	}

	return &req, nil
}

func (service *RequestService) getUploadedExamples(ctx context.Context, dto models.RequestDTOCreate) ([]types.UploadedExample, error) {
	uploadedExamples := make([]types.UploadedExample, 0)

	for i, ed := range dto.ExpectedDocuments {
		if ed.ExampleFile == nil {
			continue
		}

		if err := ValidateFileInfo(ed.ExampleFileName, ed.ExampleFileSize); err != nil {
			service.removeUploadedExamples(ctx, uploadedExamples)
			return nil, err
		}

		s3Key := service.fileStorage.GenerateS3Key(ed.ExampleFileName, "examples")
		result, err := service.fileStorage.UploadFile(ctx, s3Key, ed.ExampleFile, ed.ExampleMimeType)
		if err != nil {
			service.logger.Error("s3 upload failed for example file",
				slog.String("key", s3Key),
				slog.Any("error", err),
			)

			service.removeUploadedExamples(ctx, uploadedExamples)
			return nil, errors.ErrBadGateway{Msg: fmt.Sprintf("Failed to upload example file to S3. %v", err)}
		}

		uploadedExamples = append(uploadedExamples, types.UploadedExample{
			Index:       i,
			S3Key:       s3Key,
			S3VersionID: result.VersionId,
			MimeType:    ed.ExampleMimeType,
		})
	}

	return uploadedExamples, nil
}

func getExpectedDocuments(dto models.RequestDTOCreate, uploadedExamples []types.UploadedExample) []models.ExpectedDocument {
	expectedDocs := make([]models.ExpectedDocument, len(dto.ExpectedDocuments))
	for i, ed := range dto.ExpectedDocuments {
		expectedDocs[i] = models.ExpectedDocument{
			Title:       ed.Title,
			Description: ed.Description,
			Status:      "pending",
		}
	}
	for _, uploaded := range uploadedExamples {
		expectedDocs[uploaded.Index].ExampleFilePath = &uploaded.S3Key
		expectedDocs[uploaded.Index].ExampleMimeType = &uploaded.MimeType
	}

	return expectedDocs
}

func (service *RequestService) createRequestTransaction(ctx context.Context, req models.Request, expectedDocs []models.ExpectedDocument) (*int, error) {
	var id int
	err := service.txManager.WithTx(ctx, func(tx *sql.Tx) error {
		var txErr error
		id, txErr = service.requestRepo.AddRequestWithTx(ctx, req, tx)
		if txErr != nil {
			return txErr
		}

		for _, ed := range expectedDocs {
			ed.RequestID = id
			if txErr = service.expectedDocRepo.AddExpectedDocumentToRequestWithTx(ctx, tx, ed); txErr != nil {
				return txErr
			}
		}
		return nil
	})

	return &id, err
}

func (service *RequestService) removeUploadedExamples(ctx context.Context, uploadedExamples []types.UploadedExample) {
	for _, uploaded := range uploadedExamples {
		cleanupCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		cleanupErr := service.fileStorage.DeleteFile(cleanupCtx, uploaded.S3Key, uploaded.S3VersionID)
		cancel()
		if cleanupErr != nil {
			service.logger.Error("failed to cleanup example file after transaction failure",
				slog.String("key", uploaded.S3Key),
				slog.Any("error", cleanupErr),
			)
		}
	}
}

func ValidateFileInfo(fileName string, fileSize int64) error {
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
