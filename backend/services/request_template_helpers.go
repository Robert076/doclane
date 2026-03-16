package services

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/robfig/cron/v3"
)

func (s *RequestTemplateService) checkUserOwnsTemplate(ctx context.Context, jwtUserID int, requestTemplateID int) (*models.RequestTemplate, error) {
	requestTemplate, err := s.templateRepo.GetRequestTemplateByID(ctx, requestTemplateID)
	if err != nil {
		s.logger.Error("failed to retrieve template by id",
			slog.Int("template_id", requestTemplateID),
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err),
		)

		return nil, err
	}

	if requestTemplate.CreatedBy != jwtUserID {
		s.logger.Warn("unauthorized access attempted for a request template",
			slog.Int("template_id", requestTemplateID),
			slog.Int("user_id", jwtUserID),
		)

		return nil, errors.ErrForbidden{Msg: "This template does not belong to you."}
	}

	return &requestTemplate, nil
}

func validateRequestTemplatePatchDTO(dto models.RequestTemplateDTOPatch) error {
	if dto.Title != nil {
		if strings.TrimSpace(*dto.Title) == "" {
			return errors.ErrBadRequest{Msg: "Title cannot be empty."}
		}
		if len(*dto.Title) > 255 {
			return errors.ErrBadRequest{Msg: "Title cannot exceed 255 characters."}
		}
	}

	if dto.Description != nil && len(*dto.Description) > 1000 {
		return errors.ErrBadRequest{Msg: "Description cannot exceed 1000 characters."}
	}

	if dto.IsRecurring != nil && *dto.IsRecurring && (dto.RecurrenceCron == nil || strings.TrimSpace(*dto.RecurrenceCron) == "") {
		return errors.ErrBadRequest{Msg: "Recurrence cron is required when is_recurring is true."}
	}

	if dto.RecurrenceCron != nil && strings.TrimSpace(*dto.RecurrenceCron) != "" {
		if _, err := cron.ParseStandard(*dto.RecurrenceCron); err != nil {
			return errors.ErrBadRequest{Msg: "Invalid cron expression."}
		}
	}

	return nil
}
