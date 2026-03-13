package services

import (
	"fmt"
	"time"

	"github.com/Robert076/doclane/backend/models"
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
			if ed.Status != "uploaded" || ed.Status != "accepted" {
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
