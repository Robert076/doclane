package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/mail"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo   repositories.IUserRepo
	logger *slog.Logger
}

type RegisterParams struct {
	Email        string
	FirstName    string
	LastName     string
	Password     string
	Role         string
	DepartmentID *int
}

type LoginParams struct {
	Email    string
	Password string
}

func NewUserService(repo repositories.IUserRepo, logger *slog.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}

func (service *UserService) GetUsers(ctx context.Context, caller types.JWTClaims, limit *int, offset *int, orderBy *string, order *string, search *string) ([]models.User, error) {
	users, err := service.repo.GetUsers(ctx, limit, offset, orderBy, order, search)
	if err != nil {
		service.logger.Error("failed to fetch users",
			slog.Int("jwt_user_id", caller.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	service.logger.Info("fetched users successfully",
		slog.Int("jwt_user_id", caller.UserID),
	)
	return users, nil
}

func (service *UserService) GetUserByID(ctx context.Context, caller types.JWTClaims, id int) (*models.User, error) {
	user, err := service.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			service.logger.Warn("user not found",
				slog.Int("user_id", id),
				slog.Int("jwt_user_id", caller.UserID),
			)
			return nil, err
		}

		service.logger.Error("failed to fetch user by id",
			slog.Int("user_id", id),
			slog.Int("jwt_user_id", caller.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	service.logger.Info("fetched user by id successfully",
		slog.Int("user_id", id),
		slog.Int("jwt_user_id", caller.UserID),
	)
	return &user, nil
}

func (service *UserService) GetUserByEmail(ctx context.Context, caller types.JWTClaims, email string) (*models.User, error) {
	user, err := service.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.IsNotFound(err) {
			service.logger.Warn("user not found",
				slog.String("email", email),
				slog.Int("jwt_user_id", caller.UserID),
			)
			return nil, err
		}

		service.logger.Error("failed to fetch user by email",
			slog.String("email", email),
			slog.Int("jwt_user_id", caller.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	service.logger.Info("retrieved user by email successfully",
		slog.String("email", email),
		slog.Int("jwt_user_id", caller.UserID),
	)
	return &user, nil
}

func (service *UserService) AddUser(ctx context.Context, params RegisterParams) (int, error) {
	if err := service.ValidateUserForRegister(ctx, params.Email, params.Password, params.Role, params.FirstName, params.LastName); err != nil {
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
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		PasswordHash: string(hashedPassword),
		Role:         params.Role,
		DepartmentID: params.DepartmentID,
		IsActive:     true,
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

func (service *UserService) NotifyUser(ctx context.Context, caller types.JWTClaims, id int) error {
	user, err := service.repo.GetUserByID(ctx, id)
	if err != nil {
		service.logger.Error("error retrieving user when trying to notify",
			slog.Int("user_id", id),
			slog.Int("jwt_user_id", caller.UserID),
			slog.Any("error", err),
		)
		return err
	}

	if user.DepartmentID != nil || user.Role == types.RoleAdmin {
		service.logger.Warn("notification attempt was rejected because only normal users can be notified",
			slog.Int("user_id", id),
			slog.Int("jwt_user_id", caller.UserID),
		)
		return errors.ErrBadRequest{Msg: "This user cannot be notified"}
	}

	if !caller.IsAdmin() && !caller.IsDepartmentMember() {
		service.logger.Warn("notification attempt was rejected because user has insufficient permissions",
			slog.Int("user_id", id),
			slog.Int("jwt_user_id", caller.UserID),
		)
		return errors.ErrForbidden{Msg: "This user cannot be notified"}
	}

	if user.LastNotified != nil && user.LastNotified.After(time.Now().Add(-5*time.Minute)) {
		service.logger.Warn("notification attempt was rejected because the user was already recently notified",
			slog.Int("user_id", id),
			slog.Int("jwt_user_id", caller.UserID),
		)
		return errors.ErrTooManyRequests{Msg: fmt.Sprintf("%s %s has already been notified in the last 5 minutes.", user.FirstName, user.LastName)}
	}

	if err := service.repo.NotifyUser(ctx, id, time.Now()); err != nil {
		service.logger.Error("notification attempt failed with db error",
			slog.Int("user_id", id),
			slog.Int("jwt_user_id", caller.UserID),
			slog.Any("error", err),
		)
		return err
	}

	service.logger.Info("user notified successfully",
		slog.Int("user_id", id),
		slog.Int("jwt_user_id", caller.UserID),
	)
	return nil
}

func (service *UserService) Login(ctx context.Context, params LoginParams) (*models.User, error) {
	user, err := service.repo.GetUserByEmail(ctx, params.Email)
	if err != nil {
		if errors.IsNotFound(err) {
			service.logger.Warn("failed login attempt: user does not exist",
				slog.String("email", params.Email),
			)
			return nil, errors.ErrUnauthorized{Msg: "Invalid email or password."}
		}

		service.logger.Error("database error during login",
			slog.Any("error", err),
		)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password)); err != nil {
		service.logger.Warn("failed login attempt: incorrect password",
			slog.String("email", params.Email),
		)
		return nil, errors.ErrUnauthorized{Msg: "Invalid email or password."}
	}

	service.logger.Info("successful login",
		slog.Int("user_id", user.ID),
		slog.String("email", params.Email),
	)
	return &user, nil
}

func (service *UserService) DeactivateUser(ctx context.Context, caller types.JWTClaims, id int) error {
	_, err := service.repo.GetUserByID(ctx, id)
	if err != nil {
		service.logger.Error("could not retrieve user from db",
			slog.Int("user_id", id),
			slog.Int("jwt_user_id", caller.UserID),
			slog.Any("error", err),
		)
		return err
	}

	if !caller.IsAdmin() && !caller.IsDepartmentMember() {
		service.logger.Warn("unauthorized attempt to deactivate account",
			slog.Int("user_id", id),
			slog.Int("jwt_user_id", caller.UserID),
		)
		return errors.ErrBadRequest{Msg: "Invalid account deactivation attempted."}
	}

	if err := service.repo.DeactivateUser(ctx, id); err != nil {
		service.logger.Error("error when trying to deactivate user",
			slog.Int("user_id", id),
			slog.Int("jwt_user_id", caller.UserID),
			slog.Any("error", err),
		)
		return errors.ErrInternalServerError{Msg: fmt.Sprintf("Error deactivating user: %v", err)}
	}

	service.logger.Info("user deactivated successfully",
		slog.Int("user_id", id),
		slog.Int("jwt_user_id", caller.UserID),
	)
	return nil
}

func (service *UserService) ValidateUserForRegister(ctx context.Context, email, password, role, firstName, lastName string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.ErrBadRequest{Msg: fmt.Sprintf("Invalid email received: %s", email)}
	}

	if !types.IsValidRole(role) {
		return errors.ErrBadRequest{Msg: fmt.Sprintf("Invalid role: %s.", role)}
	}

	if firstName == "" || lastName == "" {
		return errors.ErrBadRequest{Msg: "First and last name cannot be empty."}
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
