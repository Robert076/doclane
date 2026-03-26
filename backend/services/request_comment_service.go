package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
)

type RequestCommentService struct {
	commentRepo repositories.IRequestCommentRepo
	requestRepo repositories.IRequestRepo
	userRepo    repositories.IUserRepo
	logger      *slog.Logger
}

func NewRequestCommentService(
	commentRepo repositories.IRequestCommentRepo,
	requestRepo repositories.IRequestRepo,
	userRepo repositories.IUserRepo,
	logger *slog.Logger,
) *RequestCommentService {
	return &RequestCommentService{
		commentRepo: commentRepo,
		requestRepo: requestRepo,
		userRepo:    userRepo,
		logger:      logger,
	}
}

func (s *RequestCommentService) GetCommentByID(ctx context.Context, jwtUserID int, commentID int) (*models.RequestCommentDTO, error) {
	comm, err := s.checkUserHasAccessToReadComment(ctx, jwtUserID, commentID)
	return comm, err
}

func (s *RequestCommentService) GetCommentsByRequestID(ctx context.Context, jwtUserID int, requestID int) ([]models.RequestCommentDTO, error) {
	if _, err := s.checkUserIsParticipantOfRequest(ctx, jwtUserID, requestID); err != nil {
		return nil, err
	}

	comments, err := s.commentRepo.GetCommentsByRequestID(ctx, requestID)
	return comments, err
}

func (s *RequestCommentService) AddComment(ctx context.Context, jwtUserID int, requestID int, comment models.RequestComment) (*int, error) {
	if err := s.validateComment(comment); err != nil {
		return nil, err
	}

	if err := s.checkUserIsNotSpamming(ctx, jwtUserID); err != nil {
		return nil, err
	}

	comment.UserID = jwtUserID
	comment.RequestID = requestID
	comment.CreatedAt = time.Now().UTC()
	comment.UpdatedAt = time.Now().UTC()

	if _, err := s.checkUserIsParticipantOfRequest(ctx, jwtUserID, requestID); err != nil {
		return nil, err
	}

	id, err := s.commentRepo.AddComment(ctx, comment)
	return &id, err
}
