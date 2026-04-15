package services

import (
	"context"
	"log/slog"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
)

type StatsService struct {
	db     repositories.IStatsRepo
	logger *slog.Logger
}

func NewStatsService(db repositories.IStatsRepo, logger *slog.Logger) *StatsService {
	return &StatsService{db: db, logger: logger}
}

func (s *StatsService) GetStats(ctx context.Context, claims types.JWTClaims) (*models.Stats, error) {
	if !claims.IsAdmin() {
		s.logger.Warn("unauthorized attempt to get stats",
			slog.Int("jwt_user_id", claims.UserID),
		)
		return nil, errors.ErrForbidden{Msg: "Only admins can view stats."}
	}

	stats, err := s.db.GetStats(ctx)
	if err != nil {
		s.logger.Error("failed to get stats",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	s.logger.Info("stats retrieved successfully",
		slog.Int("jwt_user_id", claims.UserID),
	)
	return stats, nil
}
