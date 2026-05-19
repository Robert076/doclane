package services

import (
	"context"
	"log/slog"

	"github.com/Robert076/doclane/backend/events"
	"github.com/Robert076/doclane/backend/repositories"
)

type AuditLogService struct {
	repo   repositories.IAuditLogRepo
	logger *slog.Logger
}

func NewAuditLogService(repo repositories.IAuditLogRepo, logger *slog.Logger) *AuditLogService {
	return &AuditLogService{repo: repo, logger: logger}
}

// OnEvent implements events.IObserver.
// It writes every event to the audit_log table.
// Errors are logged but never returned — audit logging must never
// block or fail the main request flow.
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
