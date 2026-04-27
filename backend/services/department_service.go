package services

import (
	"context"
	"log/slog"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
)

type DepartmentService struct {
	repo   repositories.IDepartmentRepo
	logger *slog.Logger
}

func NewDepartmentService(repo repositories.IDepartmentRepo, logger *slog.Logger) *DepartmentService {
	return &DepartmentService{repo: repo, logger: logger}
}

func (s *DepartmentService) GetAllDepartments(ctx context.Context, claims types.JWTClaims) ([]models.DepartmentDTORead, error) {
	departments, err := s.repo.GetAllDepartments(ctx)
	if err != nil {
		s.logger.Error("failed to fetch departments",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	s.logger.Info("fetched departments successfully",
		slog.Int("jwt_user_id", claims.UserID),
	)
	return departments, nil
}

func (s *DepartmentService) CreateDepartment(ctx context.Context, claims types.JWTClaims, name string) (int, error) {
	if !claims.IsAdmin() {
		s.logger.Warn("unauthorized attempt to create department",
			slog.Int("jwt_user_id", claims.UserID),
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
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return 0, err
	}

	s.logger.Info("department created successfully",
		slog.Int("department_id", id),
		slog.String("name", name),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return id, nil
}
