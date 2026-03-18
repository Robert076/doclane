package services

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types/errors"
)

func (s *RequestCommentService) validateComment(comment models.RequestComment) error {
	trimmed := strings.TrimSpace(comment.Comment)

	if len(trimmed) < 3 {
		return errors.ErrBadRequest{Msg: "Comment must contain at least 3 visible characters."}
	}

	if len(trimmed) > 200 {
		return errors.ErrBadRequest{Msg: "Comment is too long (max 200 characters)."}
	}

	return nil
}

func (s *RequestCommentService) checkUserIsParticipantOfRequest(ctx context.Context, jwtUserID int, requestID int) (*models.RequestDTORead, error) {
	req, err := s.requestRepo.GetRequestByID(ctx, requestID)
	if err != nil {
		s.logger.Error("error getting request from db",
			slog.Int("user_id", jwtUserID),
			slog.Int("request_id", requestID),
			slog.Any("error", err),
		)
		return nil, err
	}

	if req.ProfessionalID != jwtUserID && req.ClientID != jwtUserID {
		s.logger.Warn("unauthorized access attempted for comments on a request",
			slog.Int("user_id", jwtUserID),
			slog.Int("request_id", requestID),
		)
		return nil, errors.ErrForbidden{Msg: "You don't have access to this request."}
	}

	return &req, nil
}

func (s *RequestCommentService) checkUserHasAccessToReadComment(ctx context.Context, jwtUserID int, commentID int) (*models.RequestCommentDTO, error) {
	comm, err := s.commentRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		s.logger.Error("error getting comment from db",
			slog.Int("user_id", jwtUserID),
			slog.Int("comment_id", commentID),
			slog.Any("error", err),
		)
		return nil, err
	}

	req, err := s.requestRepo.GetRequestByID(ctx, comm.RequestID)
	if err != nil {
		s.logger.Error("error getting request from db",
			slog.Int("user_id", jwtUserID),
			slog.Int("request_id", comm.RequestID),
			slog.Any("error", err),
		)
		return nil, err
	}

	if req.ClientID != jwtUserID && req.ProfessionalID != jwtUserID {
		s.logger.Warn("unauthorized access attempt for a request comment",
			slog.Int("user_id", jwtUserID),
			slog.Int("comment_id", commentID),
		)
		return nil, errors.ErrForbidden{Msg: "Forbidden."}
	}

	return &comm, nil
}
