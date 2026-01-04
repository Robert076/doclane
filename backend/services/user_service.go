package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/mail"
	"strconv"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo   repositories.IUserRepository
	logger *slog.Logger
}

type CreateUserParams struct {
	Email          string
	Password       string
	Role           string
	ProfessionalID *int
}

func NewUserService(repo repositories.IUserRepository, logger *slog.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}

func (service *UserService) GetUsers(ctx context.Context, limit *int, offset *int, orderBy *string, order *string) ([]models.User, error) {
	users, err := service.repo.GetUsers(ctx, limit, offset, orderBy, order)
	if err != nil {
		service.logger.Error("failed to fetch users list", slog.Any("error", err))
		return nil, err
	}
	return users, nil
}

func (service *UserService) AddUser(ctx context.Context, params CreateUserParams) (int, error) {
	service.logger.Info("attempting to register new user",
		slog.String("email", params.Email),
		slog.String("role", params.Role),
	)

	if err := service.ValidateUserForRegister(ctx, params.Email, params.Password, params.Role); err != nil {
		service.logger.Warn("user validation failed for register",
			slog.String("email", params.Email),
			slog.Any("error", err),
		)
		return 0, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		service.logger.Error("failed to generate password hash",
			slog.String("email", params.Email),
			slog.Any("error", err),
		)
		return 0, err
	}

	user := models.User{
		Email:        params.Email,
		PasswordHash: string(hashedPassword),
		Role:         params.Role,
		IsActive:     true,
	}

	if params.ProfessionalID != nil {
		profIdStr := strconv.Itoa(*params.ProfessionalID)
		user.ProfessionalID = &profIdStr
	}

	id, err := service.repo.AddUser(ctx, user)
	if err != nil {
		service.logger.Error("failed to save user to database",
			slog.String("email", params.Email),
			slog.Any("error", err),
		)
		return 0, err
	}

	service.logger.Info("user registered successfully",
		slog.Int("user_id", id),
		slog.String("email", params.Email),
	)

	return id, nil
}

func (service *UserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := service.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.IsNotFound(err) {
			return models.User{}, err
		}

		service.logger.Error("failed to fetch user by email",
			slog.String("email", email),
			slog.Any("error", err),
		)
		return models.User{}, err
	}
	return user, nil
}

func (service *UserService) ValidateUserForRegister(ctx context.Context, email string, password string, role string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.ErrBadRequest{Msg: fmt.Sprintf("Invalid email received: %s", email)}
	}

	if role != "PROFESSIONAL" && role != "CLIENT" {
		return errors.ErrBadRequest{Msg: fmt.Sprintf("Invalid role: %s. Allowed: PROFESSIONAL, CLIENT", role)}
	}

	_, err := service.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return errors.ErrConflict{Msg: "User already exists."}
	}

	if !errors.IsNotFound(err) {
		service.logger.Error("database error during email availability check",
			slog.String("email", email),
			slog.Any("error", err),
		)
		return errors.ErrInternalServerError{Msg: "Failed to check if user already exists."}
	}

	return nil
}
