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
	if dto.TemplateID == 0 {
		return errors.ErrBadRequest{Msg: "A template must be provided."}
	}

	if dto.DueDate != nil && dto.DueDate.Before(time.Now()) {
		return errors.ErrBadRequest{Msg: "Due date cannot be in the past."}
	}

	if dto.IsScheduled && dto.ScheduledFor == nil {
		return errors.ErrUnprocessableContent{Msg: "A request marked as scheduled must have a scheduled_for field."}
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

func (s *RequestService) checkUserCanEditRequest(ctx context.Context, claims types.JWTClaims, requestID int) (*models.RequestDTORead, error) {
	req, err := s.requestRepo.GetRequestByID(ctx, requestID)
	if err != nil {
		s.logger.Error("error getting request from db",
			slog.Int("user_id", claims.UserID),
			slog.Int("request_id", requestID),
			slog.Any("error", err),
		)
		return nil, err
	}

	if claims.IsAdmin() {
		return &req, nil
	}

	isAllowedToEdit := claims.DepartmentID != nil && *claims.DepartmentID == req.DepartmentID
	if !isAllowedToEdit {
		return nil, errors.ErrForbidden{Msg: "You don't have access to edit the request."}
	}

	return &req, nil
}

func (s *RequestService) checkUserIsParticipantOfRequest(ctx context.Context, claims types.JWTClaims, requestID int) (*models.RequestDTORead, error) {
	req, err := s.requestRepo.GetRequestByID(ctx, requestID)
	if err != nil {
		s.logger.Error("error getting request from db",
			slog.Int("user_id", claims.UserID),
			slog.Int("request_id", requestID),
			slog.Any("error", err),
		)
		return nil, err
	}

	if claims.IsAdmin() {
		return &req, nil
	}

	isDepartmentMatch := claims.DepartmentID != nil && *claims.DepartmentID == req.DepartmentID
	isAssignee := req.Assignee == claims.UserID
	if !isDepartmentMatch && !isAssignee {
		s.logger.Warn("unauthorized access attempted for request",
			slog.Int("user_id", claims.UserID),
			slog.Int("request_id", requestID),
		)
		return nil, errors.ErrForbidden{Msg: "You don't have access to this request."}
	}

	return &req, nil
}

func (service *RequestService) checkUserHasProfileConfigured(user models.User) error {
	if user.Phone == nil {
		return errors.ErrBadRequest{Msg: "You must update your phone number first."}
	}
	if user.Street == nil {
		return errors.ErrBadRequest{Msg: "You must update the street where you live first."}
	}
	if user.Locality == nil {
		return errors.ErrBadRequest{Msg: "You must update the locality where you live first."}
	}
	return nil
}

func (service *RequestService) processRecurringRequest(ctx context.Context, req models.RequestDTORead) error {
	if req.RecurrenceCron == nil {
		return nil
	}

	expectedDocTemplates, err := service.expectedDocTmplRepo.GetByRequestTemplateID(ctx, *req.RequestTemplateID)
	if err != nil {
		return err
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

	nextDueAt := ComputeNextDueAt(nil, req.RecurrenceCron)

	newReq := models.Request{
		RequestBase: models.RequestBase{
			Assignee:       req.Assignee,
			DepartmentID:   req.DepartmentID,
			Title:          req.Title,
			Description:    req.Description,
			IsRecurring:    true,
			RecurrenceCron: req.RecurrenceCron,
			IsScheduled:    false,
			NextDueAt:      nextDueAt,
			LastUploadedAt: nil,
			DueDate:        nil,
		},
		RequestTemplateID: req.RequestTemplateID,
	}

	_, err = service.createRequestTransaction(ctx, newReq, expectedDocs)
	if err != nil {
		return err
	}

	if nextDueAt != nil {
		if err := service.requestRepo.UpdateNextDueAt(ctx, req.ID, *nextDueAt); err != nil {
			service.logger.Error("failed to update next_due_at on original request",
				slog.Int("request_id", req.ID),
				slog.Any("error", err),
			)
		}
	}

	service.logger.Info("recurring request processed successfully",
		slog.Int("original_request_id", req.ID),
		slog.Int("assignee", req.Assignee),
	)
	return nil
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
