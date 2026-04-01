package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types"
)

type RequestCommentService struct {
	commentRepo repositories.IRequestCommentRepo
	requestRepo repositories.IRequestRepo
	logger      *slog.Logger
}

func NewRequestCommentService(
	commentRepo repositories.IRequestCommentRepo,
	requestRepo repositories.IRequestRepo,
	logger *slog.Logger,
) *RequestCommentService {
	return &RequestCommentService{
		commentRepo: commentRepo,
		requestRepo: requestRepo,
		logger:      logger,
	}
}

func (s *RequestCommentService) GetCommentByID(ctx context.Context, claims types.JWTClaims, commentID int) (*models.RequestCommentDTO, error) {
	return s.checkUserHasAccessToReadComment(ctx, claims, commentID)
}

func (s *RequestCommentService) GetCommentsByRequestID(ctx context.Context, claims types.JWTClaims, requestID int) ([]models.RequestCommentDTO, error) {
	if _, err := s.checkUserIsParticipantOfRequest(ctx, claims, requestID); err != nil {
		return nil, err
	}

	comments, err := s.commentRepo.GetCommentsByRequestID(ctx, requestID)
	if err != nil {
		s.logger.Error("failed to get comments by request",
			slog.Int("request_id", requestID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}
	return comments, nil
}

func (s *RequestCommentService) AddComment(ctx context.Context, claims types.JWTClaims, requestID int, comment models.RequestComment) (*int, error) {
	if err := s.validateComment(comment); err != nil {
		return nil, err
	}

	if err := s.checkUserIsNotSpamming(ctx, claims.UserID); err != nil {
		return nil, err
	}

	if _, err := s.checkUserIsParticipantOfRequest(ctx, claims, requestID); err != nil {
		return nil, err
	}

	comment.UserID = claims.UserID
	comment.RequestID = requestID
	comment.CreatedAt = time.Now().UTC()
	comment.UpdatedAt = time.Now().UTC()

	id, err := s.commentRepo.AddComment(ctx, comment)
	if err != nil {
		s.logger.Error("failed to add comment",
			slog.Int("request_id", requestID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}
	return &id, nil
}
