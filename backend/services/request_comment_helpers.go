package services

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
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

func (s *RequestCommentService) checkUserIsParticipantOfRequest(ctx context.Context, claims types.JWTClaims, requestID int) (*models.RequestDTORead, error) {
	req, err := s.requestRepo.GetRequestByID(ctx, requestID)
	if err != nil {
		s.logger.Error("error getting request from db",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("request_id", requestID),
			slog.Any("error", err),
		)
		return nil, err
	}

	if claims.IsAdmin() {
		return &req, nil
	}

	isDepartmentMatch := claims.DepartmentID != nil && *claims.DepartmentID == req.DepartmentID
	isAssignee := req.Assignee == claims.UserID
	if !isDepartmentMatch && !isAssignee {
		s.logger.Warn("unauthorized access attempted for comments on a request",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("request_id", requestID),
		)
		return nil, errors.ErrForbidden{Msg: "You don't have access to this request."}
	}

	return &req, nil
}

func (s *RequestCommentService) checkUserHasAccessToReadComment(ctx context.Context, claims types.JWTClaims, commentID int) (*models.RequestCommentDTO, error) {
	comm, err := s.commentRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		s.logger.Error("error getting comment from db",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("comment_id", commentID),
			slog.Any("error", err),
		)
		return nil, err
	}

	if _, err := s.checkUserIsParticipantOfRequest(ctx, claims, comm.RequestID); err != nil {
		return nil, err
	}

	return &comm, nil
}

func (s *RequestCommentService) checkUserIsNotSpamming(ctx context.Context, jwtUserID int) error {
	last, err := s.commentRepo.GetLastCommentFromUser(ctx, jwtUserID)
	if err != nil {
		return nil
	}

	if time.Now().UTC().Sub(last.CreatedAt.UTC()) < 30*time.Second {
		s.logger.Warn("refused to create comment for user, was timed out",
			slog.Int("jwt_user_id", jwtUserID),
		)
		return errors.ErrTooManyRequests{Msg: "Please wait before posting another comment."}
	}

	return nil
}
