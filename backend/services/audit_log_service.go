package services

import (
	"context"
	"log/slog"

	"github.com/Robert076/doclane/backend/events"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types"
)

type AuditLogService struct {
	repo   repositories.IAuditLogRepo
	logger *slog.Logger
}

func NewAuditLogService(repo repositories.IAuditLogRepo, logger *slog.Logger) *AuditLogService {
	return &AuditLogService{repo: repo, logger: logger}
}

func (s *AuditLogService) OnEvent(ctx context.Context, event events.Event) error {
	if err := s.repo.LogEvent(ctx, event); err != nil {
		s.logger.Error("failed to write audit log entry",
			slog.String("event_type", event.Type),
			slog.Int("resource_id", event.ResourceID),
			slog.String("resource_type", event.ResourceType),
			slog.Any("error", err),
		)
	}
	return nil
}

func (s *AuditLogService) GetNotifications(ctx context.Context, caller types.CallerContext, limit int) ([]events.Event, error) {
	notifications, err := s.repo.GetNotificationsForUser(
		ctx,
		caller.UserID,
		caller.DepartmentID,
		caller.IsAdmin(),
		limit,
	)
	if err != nil {
		s.logger.Error("failed to fetch notifications",
			slog.Int("caller_id", caller.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}
	return notifications, nil
}

func (s *AuditLogService) GetByResource(ctx context.Context, resourceType string, resourceID int) ([]events.Event, error) {
	entries, err := s.repo.GetByResource(ctx, resourceType, resourceID)
	if err != nil {
		s.logger.Error("failed to fetch audit log",
			slog.String("resource_type", resourceType),
			slog.Int("resource_id", resourceID),
			slog.Any("error", err),
		)
		return nil, err
	}

	return entries, nil
}
