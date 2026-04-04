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
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
)

type InvitationCodeService struct {
	invitationRepo repositories.IInvitationCodeRepo
	departmentRepo repositories.IDepartmentRepo
	logger         *slog.Logger
}

func NewInvitationCodeService(
	invitationRepo repositories.IInvitationCodeRepo,
	departmentRepo repositories.IDepartmentRepo,
	logger *slog.Logger,
) *InvitationCodeService {
	return &InvitationCodeService{
		invitationRepo: invitationRepo,
		departmentRepo: departmentRepo,
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
	claims types.JWTClaims,
	departmentID int,
	expiresInDays int,
) (string, error) {
	if !claims.IsAdmin() {
		s.logger.Warn("non-admin attempted to create invitation code",
			slog.Int("jwt_user_id", claims.UserID),
		)
		return "", errors.ErrForbidden{Msg: "Only admins can generate invitation codes."}
	}

	if _, err := s.departmentRepo.GetDepartmentByID(ctx, departmentID); err != nil {
		s.logger.Warn("admin tried to create invitation code for department that does not exist",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("department_id", departmentID),
		)
		return "", errors.ErrNotFound{Msg: "Department not found."}
	}

	existingCodes, err := s.invitationRepo.GetInvitationCodesByCreator(ctx, claims.UserID)
	if err != nil {
		s.logger.Error("error getting codes for user",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return "", err
	}

	if len(existingCodes) >= 3 {
		s.logger.Warn("user tried to generate too many active invitation codes",
			slog.Int("jwt_user_id", claims.UserID),
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

	if err = s.invitationRepo.CreateInvitationCode(ctx, departmentID, code, claims.UserID, expiresAt); err != nil {
		s.logger.Error("failed to save invitation code to database",
			slog.String("code", code),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return "", errors.ErrInternalServerError{Msg: "Failed to save invitation code."}
	}

	s.logger.Info("invitation code created successfully",
		slog.String("code", code),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return code, nil
}

func (s *InvitationCodeService) GetInvitationCodes(
	ctx context.Context,
	claims types.JWTClaims,
) ([]models.InvitationCode, error) {
	if !claims.IsAdmin() {
		return nil, errors.ErrForbidden{Msg: "Only admins can view invitation codes."}
	}

	codes, err := s.invitationRepo.GetInvitationCodesByCreator(ctx, claims.UserID)
	if err != nil {
		s.logger.Error("failed to fetch invitation codes",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	validCodes, err := s.deleteExpiredCodes(ctx, codes)
	if err != nil {
		return nil, err
	}

	s.logger.Info("invitation codes retrieved successfully",
		slog.Int("jwt_user_id", claims.UserID),
		slog.Int("count", len(validCodes)),
	)
	return validCodes, nil
}

func (s *InvitationCodeService) GetInvitationCodesByDepartment(
	ctx context.Context,
	claims types.JWTClaims,
	departmentID int,
) ([]models.InvitationCode, error) {
	if !claims.IsAdmin() {
		s.logger.Warn("non-admin attempted to view invitation codes by department",
			slog.Int("jwt_user_id", claims.UserID),
		)
		return nil, errors.ErrForbidden{Msg: "Only admins can view invitation codes."}
	}

	if _, err := s.departmentRepo.GetDepartmentByID(ctx, departmentID); err != nil {
		s.logger.Warn("tried to get invitation codes for department that does not exist",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("department_id", departmentID),
		)
		return nil, errors.ErrNotFound{Msg: "Department not found."}
	}

	codes, err := s.invitationRepo.GetInvitationCodesByDepartment(ctx, departmentID)
	if err != nil {
		s.logger.Error("failed to fetch invitation codes by department",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("department_id", departmentID),
			slog.Any("error", err),
		)
		return nil, err
	}

	validCodes, err := s.deleteExpiredCodes(ctx, codes)
	if err != nil {
		return nil, err
	}

	s.logger.Info("invitation codes by department retrieved successfully",
		slog.Int("jwt_user_id", claims.UserID),
		slog.Int("department_id", departmentID),
		slog.Int("count", len(validCodes)),
	)
	return validCodes, nil
}

func (s *InvitationCodeService) ValidateAndUseInvitationCode(
	ctx context.Context,
	code string,
) error {
	invCode, err := s.invitationRepo.GetInvitationCodeByCode(ctx, code)
	if err != nil {
		s.logger.Warn("invitation code not found",
			slog.String("code", code),
			slog.Any("error", err),
		)
		return errors.ErrNotFound{Msg: "Invalid invitation code."}
	}

	if invCode.UsedAt != nil {
		s.logger.Warn("attempted to use already-used invitation code",
			slog.String("code", code),
		)
		return errors.ErrBadRequest{Msg: "This invitation code has already been used."}
	}

	if invCode.ExpiresAt != nil && time.Now().After(*invCode.ExpiresAt) {
		s.logger.Warn("attempted to use expired invitation code",
			slog.String("code", code),
			slog.Time("expired_at", *invCode.ExpiresAt),
		)
		return errors.ErrBadRequest{Msg: "This invitation code has expired."}
	}

	if err = s.invitationRepo.InvalidateCode(ctx, invCode.ID); err != nil {
		s.logger.Error("failed to invalidate invitation code",
			slog.Int("code_id", invCode.ID),
			slog.Any("error", err),
		)
		return errors.ErrInternalServerError{Msg: "Failed to process invitation code."}
	}

	s.logger.Info("invitation code used successfully",
		slog.String("code", code),
		slog.Int("created_by", invCode.CreatedBy),
	)
	return nil
}

func (s *InvitationCodeService) DeleteInvitationCode(
	ctx context.Context,
	claims types.JWTClaims,
	codeID int,
) error {
	if !claims.IsAdmin() {
		return errors.ErrForbidden{Msg: "Only admins can delete invitation codes."}
	}

	if _, err := s.checkUserOwnsInvitationCode(ctx, claims.UserID, codeID); err != nil {
		return err
	}

	if err := s.invitationRepo.DeleteCode(ctx, codeID); err != nil {
		s.logger.Error("failed to delete invitation code",
			slog.Int("code_id", codeID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return errors.ErrInternalServerError{Msg: "Failed to delete invitation code."}
	}

	s.logger.Info("invitation code deleted successfully",
		slog.Int("code_id", codeID),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return nil
}
