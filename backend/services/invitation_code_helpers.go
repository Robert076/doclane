package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types/errors"
)

func (s *InvitationCodeService) checkUserOwnsInvitationCode(ctx context.Context, jwtUserID int, invitationCodeID int) (*models.InvitationCode, error) {
	code, err := s.invitationRepo.GetInvitationCodeByID(ctx, invitationCodeID)
	if err != nil {
		s.logger.Warn("invitation code not found",
			slog.Int("code_id", invitationCodeID),
			slog.Any("error", err),
		)
		return nil, errors.ErrNotFound{Msg: "Invitation code not found."}
	}

	if code.CreatedBy != jwtUserID {
		s.logger.Warn("user attempted to access another user's invitation code",
			slog.Int("jwt_user_id", jwtUserID),
			slog.Int("code_created_by", code.CreatedBy),
			slog.Int("code_id", invitationCodeID),
		)
		return nil, errors.ErrForbidden{Msg: "You can only access your own invitation codes."}
	}

	return &code, nil
}

func (s *InvitationCodeService) deleteExpiredCodes(ctx context.Context, codes []models.InvitationCode) ([]models.InvitationCode, error) {
	validCodes := make([]models.InvitationCode, 0)
	for _, code := range codes {
		if code.ExpiresAt != nil && code.ExpiresAt.Before(time.Now()) {
			if err := s.invitationRepo.DeleteCode(ctx, code.ID); err != nil {
				s.logger.Error("failed to delete expired code",
					slog.Int("code_id", code.ID),
					slog.Int("created_by", code.CreatedBy),
					slog.Any("error", err),
				)
				return nil, errors.ErrInternalServerError{Msg: fmt.Sprintf("Could not delete expired token. %v", err)}
			}
		} else {
			validCodes = append(validCodes, code)
		}
	}
	return validCodes, nil
}
