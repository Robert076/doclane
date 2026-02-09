package services

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

func ValidateRequestInput(dto models.DocumentRequestDTOCreate) error {
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

	return nil
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

func ComputeStatus(lastUploadedAt *time.Time, nextDueAt *time.Time) string {
	now := time.Now()

	if nextDueAt == nil && lastUploadedAt != nil {
		return "uploaded"
	}

	if nextDueAt == nil {
		return "pending"
	}

	if now.After(*nextDueAt) && lastUploadedAt == nil {
		return "overdue"
	}

	if lastUploadedAt != nil {
		return "uploaded"
	}

	return "pending"
}

func generateS3Key(fileName string, requestID int) string {
	cleanFileName := filepath.Base(fileName)
	uniqueID := uuid.New().String()

	s3Key := fmt.Sprintf("requests/%d/%s-%s", requestID, uniqueID, cleanFileName)

	return s3Key
}

func ValidatePatchDTO(dto models.DocumentRequestDTOPatch) error {
	if len(dto.Title) < 3 || len(dto.Title) > 30 {
		return errors.ErrBadRequest{Msg: "New title is too short or too long. Minimum 3 characters, maximum 30 characters."}
	}

	return nil
}
