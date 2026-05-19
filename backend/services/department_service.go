package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/Robert076/doclane/backend/events"
	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
)

type DepartmentService struct {
	repo   repositories.IDepartmentRepo
	logger *slog.Logger
	bus    *events.EventBus
}

func NewDepartmentService(repo repositories.IDepartmentRepo, logger *slog.Logger, bus *events.EventBus) *DepartmentService {
	return &DepartmentService{repo: repo, logger: logger, bus: bus}
}

func (s *DepartmentService) GetAllDepartments(ctx context.Context, claims types.CallerContext) ([]models.DepartmentDTORead, error) {
	departments, err := s.repo.GetAllDepartments(ctx)
	if err != nil {
		s.logger.Error("failed to fetch departments",
			slog.Int("caller_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	s.logger.Info("fetched departments successfully",
		slog.Int("caller_id", claims.UserID),
	)
	return departments, nil
}

func (s *DepartmentService) CreateDepartment(ctx context.Context, claims types.CallerContext, name string) (int, error) {
	if !claims.IsAdmin() {
		s.logger.Warn("unauthorized attempt to create department",
			slog.Int("caller_id", claims.UserID),
		)
		return 0, errors.ErrForbidden{Msg: "Only admins can create departments."}
	}

	if name == "" {
		return 0, errors.ErrBadRequest{Msg: "Department name cannot be empty."}
	}

	id, err := s.repo.CreateDepartment(ctx, name)
	if err != nil {
		s.logger.Error("failed to create department",
			slog.String("name", name),
			slog.Int("caller_id", claims.UserID),
			slog.Any("error", err),
		)
		return 0, err
	}

	s.logger.Info("department created successfully",
		slog.Int("department_id", id),
		slog.String("name", name),
		slog.Int("caller_id", claims.UserID),
	)

	s.bus.Publish(ctx, events.Event{
		Type:         events.EventDepartmentCreated,
		ActorID:      claims.UserID,
		ResourceID:   id,
		ResourceType: events.ResourceTypeDepartment,
		Metadata: map[string]any{
			"name": name,
		},
		OccurredAt: time.Now().UTC(),
	})
	return id, nil
}
