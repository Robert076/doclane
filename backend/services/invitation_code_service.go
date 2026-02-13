// services/invitation_code_service.go
package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types/errors"
)

type InvitationCodeService struct {
	invitationRepo repositories.IInvitationCodeRepository
	userRepo       repositories.IUserRepository
	logger         *slog.Logger
}

func NewInvitationCodeService(
	invitationRepo repositories.IInvitationCodeRepository,
	userRepo repositories.IUserRepository,
	logger *slog.Logger,
) *InvitationCodeService {
	return &InvitationCodeService{
		invitationRepo: invitationRepo,
		userRepo:       userRepo,
		logger:         logger,
	}
}

func generateInvitationCode() (string, error) {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	code := strings.ToUpper(hex.EncodeToString(bytes))
	return fmt.Sprintf("%s-%s-%s", code[0:4], code[4:8], code[8:12]), nil
}

func (s *InvitationCodeService) CreateInvitationCode(
	ctx context.Context,
	jwtUserId int,
	expiresInDays int,
) (string, error) {
	user, err := s.userRepo.GetUserByID(ctx, jwtUserId)
	if err != nil {
		s.logger.Error("failed to fetch user for invitation code creation",
			slog.Int("user_id", jwtUserId),
			slog.Any("error", err),
		)
		return "", err
	}

	if user.Role != "PROFESSIONAL" {
		s.logger.Warn("non-professional attempted to create invitation code",
			slog.Int("user_id", jwtUserId),
			slog.String("role", user.Role),
		)
		return "", errors.ErrForbidden{Msg: "Only professionals can generate invitation codes."}
	}

	existingCodes, err := s.invitationRepo.GetInvitationCodesByProfessional(ctx, jwtUserId)
	if err != nil {
		s.logger.Error("error getting codes for professional",
			slog.Int("user_id", jwtUserId),
			slog.Any("error", err),
		)
	}

	if len(existingCodes) >= 3 {
		s.logger.Warn("user tried to generate too many active invitation codes",
			slog.Int("user_id", jwtUserId),
			slog.Any("error", err),
		)
		return "", errors.ErrBadRequest{Msg: "Only 3 active codes are allowed at one time."}
	}

	code, err := generateInvitationCode()
	if err != nil {
		s.logger.Error("failed to generate invitation code",
			slog.Any("error", err),
		)
		return "", errors.ErrInternalServerError{Msg: "Failed to generate invitation code."}
	}

	var expiresAt *time.Time
	if expiresInDays > 0 {
		expiry := time.Now().Add(time.Duration(expiresInDays) * 24 * time.Hour)
		expiresAt = &expiry
	}

	err = s.invitationRepo.CreateInvitationCode(ctx, code, jwtUserId, expiresAt)
	if err != nil {
		s.logger.Error("failed to save invitation code to database",
			slog.String("code", code),
			slog.Int("professional_id", jwtUserId),
			slog.Any("error", err),
		)
		return "", errors.ErrInternalServerError{Msg: "Failed to save invitation code."}
	}

	s.logger.Info("invitation code created successfully",
		slog.String("code", code),
		slog.Int("professional_id", jwtUserId),
	)

	return code, nil
}

func (s *InvitationCodeService) GetInvitationCodesByProfessional(
	ctx context.Context,
	jwtUserId int,
) ([]models.InvitationCode, error) {
	user, err := s.userRepo.GetUserByID(ctx, jwtUserId)
	if err != nil {
		s.logger.Error("failed to fetch user for invitation codes list",
			slog.Int("user_id", jwtUserId),
			slog.Any("error", err),
		)
		return nil, err
	}

	if user.Role != "PROFESSIONAL" {
		s.logger.Warn("non-professional attempted to view invitation codes",
			slog.Int("user_id", jwtUserId),
			slog.String("role", user.Role),
		)
		return nil, errors.ErrForbidden{Msg: "Only professionals can view invitation codes."}
	}

	codes, err := s.invitationRepo.GetInvitationCodesByProfessional(ctx, jwtUserId)
	if err != nil {
		s.logger.Error("failed to fetch invitation codes from repo",
			slog.Int("professional_id", jwtUserId),
			slog.Any("error", err),
		)
		return nil, errors.ErrInternalServerError{Msg: fmt.Sprintf("Could not fetch codes. %v", err)}
	}

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
		}
	}

	s.logger.Info("invitation codes retrieved successfully",
		slog.Int("professional_id", jwtUserId),
		slog.Int("count", len(codes)),
	)

	return codes, nil
}

func (s *InvitationCodeService) ValidateAndUseInvitationCode(
	ctx context.Context,
	code string,
) (int, error) {
	invCode, err := s.invitationRepo.GetInvitationCodeByCode(ctx, code)
	if err != nil {
		s.logger.Warn("invitation code not found",
			slog.String("code", code),
			slog.Any("error", err),
		)
		return 0, errors.ErrNotFound{Msg: "Invalid invitation code."}
	}

	if invCode.UsedAt != nil {
		s.logger.Warn("attempted to use already-used invitation code",
			slog.String("code", code),
		)
		return 0, errors.ErrBadRequest{Msg: "This invitation code has already been used."}
	}

	if invCode.ExpiresAt != nil && time.Now().After(*invCode.ExpiresAt) {
		s.logger.Warn("attempted to use expired invitation code",
			slog.String("code", code),
			slog.Time("expired_at", *invCode.ExpiresAt),
		)
		return 0, errors.ErrBadRequest{Msg: "This invitation code has expired."}
	}

	err = s.invitationRepo.InvalidateCode(ctx, invCode.ID)
	if err != nil {
		s.logger.Error("failed to invalidate invitation code",
			slog.Int("code_id", invCode.ID),
			slog.Any("error", err),
		)
		return 0, errors.ErrInternalServerError{Msg: "Failed to process invitation code."}
	}

	s.logger.Info("invitation code used successfully",
		slog.String("code", code),
		slog.Int("professional_id", invCode.ProfessionalID),
	)

	return invCode.ProfessionalID, nil
}

func (s *InvitationCodeService) GetInvitationCodeByCode(
	ctx context.Context,
	code string,
) (models.InvitationCode, error) {
	invCode, err := s.invitationRepo.GetInvitationCodeByCode(ctx, code)
	if err != nil {
		s.logger.Warn("invitation code lookup failed",
			slog.String("code", code),
			slog.Any("error", err),
		)
		return models.InvitationCode{}, errors.ErrNotFound{Msg: "Invalid invitation code."}
	}

	if invCode.UsedAt != nil {
		return models.InvitationCode{}, errors.ErrBadRequest{Msg: "This invitation code has already been used."}
	}

	if invCode.ExpiresAt != nil && time.Now().After(*invCode.ExpiresAt) {
		return models.InvitationCode{}, errors.ErrBadRequest{Msg: "This invitation code has expired."}
	}

	return invCode, nil
}

func (s *InvitationCodeService) DeleteInvitationCode(
	ctx context.Context,
	jwtUserId int,
	codeID int,
) error {
	user, err := s.userRepo.GetUserByID(ctx, jwtUserId)
	if err != nil {
		s.logger.Error("failed to fetch user for code deletion",
			slog.Int("user_id", jwtUserId),
			slog.Any("error", err),
		)
		return err
	}

	if user.Role != "PROFESSIONAL" {
		s.logger.Warn("non-professional attempted to delete invitation code",
			slog.Int("user_id", jwtUserId),
			slog.String("role", user.Role),
		)
		return errors.ErrForbidden{Msg: "Only professionals can delete invitation codes."}
	}

	code, err := s.invitationRepo.GetInvitationCodeByID(ctx, codeID)
	if err != nil {
		s.logger.Warn("invitation code not found for deletion",
			slog.Int("code_id", codeID),
			slog.Any("error", err),
		)
		return errors.ErrNotFound{Msg: "Invitation code not found."}
	}

	if code.ProfessionalID != jwtUserId {
		s.logger.Warn("professional attempted to delete another professional's code",
			slog.Int("user_id", jwtUserId),
			slog.Int("code_professional_id", code.ProfessionalID),
			slog.Int("code_id", codeID),
		)
		return errors.ErrForbidden{Msg: "You can only delete your own invitation codes."}
	}

	err = s.invitationRepo.InvalidateCode(ctx, codeID)
	if err != nil {
		s.logger.Error("failed to delete invitation code",
			slog.Int("code_id", codeID),
			slog.Any("error", err),
		)
		return errors.ErrInternalServerError{Msg: "Failed to delete invitation code."}
	}

	s.logger.Info("invitation code deleted successfully",
		slog.Int("code_id", codeID),
		slog.Int("professional_id", jwtUserId),
	)

	return nil
}

func (s *InvitationCodeService) ReactivateCode(ctx context.Context, code string) error {
	err := s.invitationRepo.ReactivateCode(ctx, code)
	if err != nil {
		s.logger.Error("failed to reactivate invitation code",
			slog.String("code", code),
			slog.Any("error", err),
		)
		return errors.ErrInternalServerError{Msg: "Failed to reactivate invitation code."}
	}

	s.logger.Warn("invitation code reactivated due to failed registration",
		slog.String("code", code),
	)

	return nil
}
