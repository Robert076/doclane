package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/mail"
	"time"

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
	FirstName      string
	LastName       string
	Password       string
	Role           string
	ProfessionalID *int
}

type LoginParams struct {
	Email    string
	Password string
}

func NewUserService(repo repositories.IUserRepository, logger *slog.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}

func (service *UserService) GetUsers(ctx context.Context, limit *int, offset *int, orderBy *string, order *string, search *string) ([]models.User, error) {
	users, err := service.repo.GetUsers(ctx, limit, offset, orderBy, order, search)
	if err != nil {
		service.logger.Error("failed to fetch users list", slog.Any("error", err))
		return nil, err
	}
	return users, nil
}

func (service *UserService) GetUserByID(ctx context.Context, id int) (models.User, error) {
	service.logger.Info("fetching user by id", slog.Int("user_id", id))

	user, err := service.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			service.logger.Warn("user not found", slog.Int("user_id", id))
			return models.User{}, errors.ErrNotFound{Msg: "User not found"}
		}

		service.logger.Error("failed to fetch user by id",
			slog.Int("user_id", id),
			slog.Any("error", err),
		)
		return models.User{}, err
	}

	return user, nil
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

func (service *UserService) GetProfessionalClients(
	ctx context.Context,
	jwtUserId int,
	limit *int,
	offset *int,
) ([]models.User, error) {
	service.logger.Info("fetching clients for professional", slog.Int("professional_id", jwtUserId))

	user, err := service.repo.GetUserByID(ctx, jwtUserId)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound{Msg: "Professional user not found."}
		}
		service.logger.Error("failed to verify professional role",
			slog.Int("user_id", jwtUserId),
			slog.Any("error", err),
		)
		return nil, err
	}

	if user.Role != "PROFESSIONAL" {
		service.logger.Warn("unauthorized access attempt to professional clients list",
			slog.Int("user_id", jwtUserId),
			slog.String("role", user.Role),
		)
		return nil, errors.ErrForbidden{Msg: "Only professionals can access this client list."}
	}

	clients, err := service.repo.GetUsersByProfessionalID(ctx, jwtUserId, limit, offset)
	if err != nil {
		service.logger.Error("failed to fetch clients for professional",
			slog.Int("professional_id", jwtUserId),
			slog.Any("error", err),
		)
		return nil, err
	}

	service.logger.Info("successfully retrieved clients",
		slog.Int("professional_id", jwtUserId),
		slog.Int("count", len(clients)),
	)

	return clients, nil
}

func (service *UserService) AddUser(ctx context.Context, params CreateUserParams) (int, error) {
	service.logger.Info("attempting to register new user",
		slog.String("email", params.Email),
		slog.String("role", params.Role),
	)

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
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if params.ProfessionalID != nil {
		user.ProfessionalID = params.ProfessionalID
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

func (service *UserService) NotifyUser(ctx context.Context, jwtUserID int, id int) error {
	user, err := service.repo.GetUserByID(ctx, id)
	if err != nil {
		service.logger.Error("error retrieving user when trying to notify",
			slog.Int("user_id", id),
			slog.Any("error", err),
		)
		return errors.ErrInternalServerError{Msg: "Could not retrieve user"}
	}

	if user.ProfessionalID == nil {
		service.logger.Warn("notification attempt was rejected because user is not client",
			slog.Int("user_id", id),
			slog.Int("jwt_user_id", jwtUserID),
		)
		return errors.ErrBadRequest{Msg: "This user cannot be notified"}
	}

	if *user.ProfessionalID != jwtUserID {
		service.logger.Warn("notification attempt was rejected due to insufficient permissions",
			slog.Int("user_id", id),
			slog.Int("jwt_user_id", jwtUserID),
		)
		return errors.ErrForbidden{Msg: "You are not allowed to notify this user."}
	}

	if user.LastNotified != nil {
		if (*user.LastNotified).After(time.Now().Add(time.Minute * -5)) {
			return errors.ErrTooManyRequests{Msg: fmt.Sprintf("%s %s has already been notified in the last 5 minutes.", user.FirstName, user.LastName)}
		}
	}

	if err := service.repo.NotifyUser(ctx, id, time.Now()); err != nil {
		service.logger.Error("notification attempt failed with db error",
			slog.Int("user_id", id),
			slog.Any("error", err),
		)

		return errors.ErrInternalServerError{Msg: "Something went wrong"}
	}

	return nil
}

func (service *UserService) Login(ctx context.Context, params LoginParams) (models.User, error) {
	service.logger.Info("login attempt", slog.String("email", params.Email))

	user, err := service.repo.GetUserByEmail(ctx, params.Email)
	if err != nil {
		if errors.IsNotFound(err) {
			return models.User{}, errors.ErrUnauthorized{Msg: "Invalid email or password."}
		}
		service.logger.Error("database error during login", slog.Any("error", err))
		return models.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password)); err != nil {
		service.logger.Warn("failed login attempt: incorrect password", slog.String("email", params.Email))
		return models.User{}, errors.ErrUnauthorized{Msg: "Invalid email or password."}
	}

	service.logger.Info("successful login", slog.Int("user_id", user.ID), slog.String("email", params.Email))
	return user, nil
}

func (service *UserService) DeactivateUser(ctx context.Context, jwtUserID int, id int) error {
	u, err := service.repo.GetUserByID(ctx, id)
	if err != nil {
		return errors.ErrBadRequest{Msg: "Invalid user received."}
	}

	if u.ProfessionalID == nil {
		return errors.ErrBadRequest{Msg: "Invalid account deactivation attempted."}
	}

	if jwtUserID != *u.ProfessionalID {
		return errors.ErrForbidden{Msg: "You cannot deactivate an account of someone who is not your client."}
	}

	if err := service.repo.DeactivateUser(ctx, id); err != nil {
		return errors.ErrInternalServerError{Msg: fmt.Sprintf("Error deactivating user: %v", err)}
	}

	return nil
}

func (service *UserService) ValidateUserForRegister(ctx context.Context, email, password, role, firstName, lastName string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.ErrBadRequest{Msg: fmt.Sprintf("Invalid email received: %s", email)}
	}

	if role != "PROFESSIONAL" && role != "CLIENT" {
		return errors.ErrBadRequest{Msg: fmt.Sprintf("Invalid role: %s. Allowed: PROFESSIONAL, CLIENT", role)}
	}

	if firstName == "" || lastName == "" {
		return errors.ErrBadRequest{Msg: fmt.Sprintf("First and last name cannot be empty.")}
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
