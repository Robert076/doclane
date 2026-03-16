package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types/errors"
)

func (s *InvitationCodeService) checkUserIsProfessional(ctx context.Context, jwtUserID int) error {
	user, err := s.userRepo.GetUserByID(ctx, jwtUserID)
	if err != nil {
		s.logger.Error("failed to fetch user",
			slog.Int("user_id", jwtUserID),
			slog.Any("error", err),
		)
		return err
	}
	if user.Role != "PROFESSIONAL" {
		s.logger.Warn("non-professional attempted to access invitation codes",
			slog.Int("user_id", jwtUserID),
			slog.String("role", user.Role),
		)
		return errors.ErrForbidden{Msg: "You are not allowed to view invitation codes."}
	}

	return nil
}

func (s *InvitationCodeService) checkUserOwnsInvitationCode(ctx context.Context, jwtUserID int, invitationCodeID int) (*models.InvitationCode, error) {
	code, err := s.invitationRepo.GetInvitationCodeByID(ctx, invitationCodeID)
	if err != nil {
		s.logger.Warn("invitation code not found",
			slog.Int("code_id", invitationCodeID),
			slog.Any("error", err),
		)
		return nil, errors.ErrNotFound{Msg: "Invitation code not found."}
	}

	if code.ProfessionalID != jwtUserID {
		s.logger.Warn("professional attempted to access another professional's code",
			slog.Int("user_id", jwtUserID),
			slog.Int("code_professional_id", code.ProfessionalID),
			slog.Int("code_id", invitationCodeID),
		)
		return nil, errors.ErrForbidden{Msg: "You can only access your own invitation codes."}
	}

	return &code, nil
}

func (s *InvitationCodeService) deleteExpiredCodes(ctx context.Context, codes []models.InvitationCode) ([]models.InvitationCode, error) {
	validCodes := make([]models.InvitationCode, 0)
	for _, code := range codes {
		if code.ExpiresAt.Before(time.Now()) {
			err := s.invitationRepo.DeleteCode(ctx, code.ID)
			if err != nil {
				s.logger.Error("failed to delete expired code",
					slog.Int("code_id", code.ID),
					slog.Int("professional_id", code.ProfessionalID),
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
