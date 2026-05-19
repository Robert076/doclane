package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Robert076/doclane/backend/events"
	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
)

type UserService struct {
	repo        repositories.IUserRepo
	requestRepo repositories.IRequestRepo
	logger      *slog.Logger
	bus         *events.EventBus
}

type SyncUserParams struct {
	CognitoSub   string
	Email        string
	FirstName    string
	LastName     string
	Role         string
	DepartmentID *int
}

func NewUserService(repo repositories.IUserRepo, requestRepo repositories.IRequestRepo, logger *slog.Logger, bus *events.EventBus) *UserService {
	return &UserService{repo: repo, requestRepo: requestRepo, logger: logger, bus: bus}
}

// SyncUser creates a DB record for a user who has already been authenticated
// by Cognito. This is called once after the user confirms their email.
// Passwords are never handled here — Cognito owns credentials.
func (service *UserService) SyncUser(ctx context.Context, params SyncUserParams) (int, error) {
	if params.FirstName == "" || params.LastName == "" {
		return 0, errors.ErrBadRequest{Msg: "First and last name cannot be empty."}
	}

	if !types.IsValidRole(params.Role) {
		return 0, errors.ErrBadRequest{Msg: fmt.Sprintf("Invalid role: %s.", params.Role)}
	}

	_, err := service.repo.GetUserByCognitoSub(ctx, params.CognitoSub)
	if err == nil {
		return 0, errors.ErrConflict{Msg: "User already exists."}
	}
	if !errors.IsNotFound(err) {
		service.logger.Error("database error during cognito_sub availability check",
			slog.String("cognito_sub", params.CognitoSub),
			slog.Any("error", err),
		)
		return 0, errors.ErrInternalServerError{Msg: "Failed to check if user already exists."}
	}

	user := models.User{
		CognitoSub:   params.CognitoSub,
		Email:        params.Email,
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Role:         params.Role,
		DepartmentID: params.DepartmentID,
		IsActive:     true,
	}

	id, err := service.repo.AddUser(ctx, user)
	if err != nil {
		service.logger.Error("failed to save synced user to database",
			slog.String("cognito_sub", params.CognitoSub),
			slog.String("email", params.Email),
			slog.Any("error", err),
		)
		return 0, err
	}

	service.logger.Info("user synced from Cognito successfully",
		slog.Int("user_id", id),
		slog.String("cognito_sub", params.CognitoSub),
		slog.String("email", params.Email),
	)
	return id, nil
}

func (service *UserService) GetUsers(ctx context.Context, caller types.CallerContext, limit *int, offset *int, orderBy *string, order *string, search *string) ([]models.User, error) {
	users, err := service.repo.GetUsers(ctx, limit, offset, orderBy, order, search)
	if err != nil {
		service.logger.Error("failed to fetch users",
			slog.Int("caller_id", caller.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	service.logger.Info("fetched users successfully",
		slog.Int("caller_id", caller.UserID),
	)
	return users, nil
}

func (service *UserService) GetUserByID(ctx context.Context, caller types.CallerContext, id int) (*models.User, error) {
	user, err := service.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			service.logger.Warn("user not found",
				slog.Int("user_id", id),
				slog.Int("caller_id", caller.UserID),
			)
			return nil, err
		}

		service.logger.Error("failed to fetch user by id",
			slog.Int("user_id", id),
			slog.Int("caller_id", caller.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	service.logger.Info("fetched user by id successfully",
		slog.Int("user_id", id),
		slog.Int("caller_id", caller.UserID),
	)
	return &user, nil
}

func (service *UserService) GetUserByEmail(ctx context.Context, caller types.CallerContext, email string) (*models.User, error) {
	user, err := service.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.IsNotFound(err) {
			service.logger.Warn("user not found",
				slog.String("email", email),
				slog.Int("caller_id", caller.UserID),
			)
			return nil, err
		}

		service.logger.Error("failed to fetch user by email",
			slog.String("email", email),
			slog.Int("caller_id", caller.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	service.logger.Info("retrieved user by email successfully",
		slog.String("email", email),
		slog.Int("caller_id", caller.UserID),
	)
	return &user, nil
}

func (service *UserService) GetUsersByDepartment(ctx context.Context, caller types.CallerContext, departmentID int) ([]models.User, error) {
	if !caller.IsAdmin() {
		service.logger.Warn("unauthorized attempt to retrieve users by department",
			slog.Int("caller_id", caller.UserID),
		)
		return nil, errors.ErrForbidden{Msg: "Only admins can list users by department."}
	}

	users, err := service.repo.GetUsersByDepartment(ctx, departmentID)
	if err != nil {
		service.logger.Error("error when trying to retrieve users by department",
			slog.Int("caller_id", caller.UserID),
			slog.Any("error", err),
		)
		return nil, err
	}

	service.logger.Info("retrieved users by department successfully",
		slog.Int("caller_id", caller.UserID),
		slog.Int("department_id", departmentID),
	)
	return users, nil
}

func (service *UserService) NotifyUser(ctx context.Context, caller types.CallerContext, id int) error {
	user, err := service.repo.GetUserByID(ctx, id)
	if err != nil {
		service.logger.Error("error retrieving user when trying to notify",
			slog.Int("user_id", id),
			slog.Int("caller_id", caller.UserID),
			slog.Any("error", err),
		)
		return err
	}

	if user.DepartmentID != nil || user.Role == types.RoleAdmin {
		service.logger.Warn("notification rejected: only regular users (citizens) can be notified",
			slog.Int("user_id", id),
			slog.Int("caller_id", caller.UserID),
		)
		return errors.ErrBadRequest{Msg: "This user cannot be notified."}
	}

	if !caller.IsAdmin() && !caller.IsDepartmentMember() {
		service.logger.Warn("notification rejected: caller has insufficient permissions",
			slog.Int("user_id", id),
			slog.Int("caller_id", caller.UserID),
		)
		return errors.ErrForbidden{Msg: "You do not have permission to notify this user."}
	}

	if user.LastNotified != nil && user.LastNotified.After(time.Now().Add(-5*time.Minute)) {
		service.logger.Warn("notification rejected: user was already recently notified",
			slog.Int("user_id", id),
			slog.Int("caller_id", caller.UserID),
		)
		return errors.ErrTooManyRequests{Msg: fmt.Sprintf("%s %s has already been notified in the last 5 minutes.", user.FirstName, user.LastName)}
	}

	if err := service.repo.NotifyUser(ctx, id, time.Now()); err != nil {
		service.logger.Error("notification failed at db layer",
			slog.Int("user_id", id),
			slog.Int("caller_id", caller.UserID),
			slog.Any("error", err),
		)
		return err
	}

	service.logger.Info("user notified successfully",
		slog.Int("user_id", id),
		slog.Int("caller_id", caller.UserID),
	)

	service.bus.Publish(ctx, events.Event{
		Type:         events.EventUserNotified,
		ActorID:      caller.UserID,
		ResourceID:   id,
		ResourceType: events.ResourceTypeUser,
		Metadata: map[string]any{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
		},
		OccurredAt: time.Now().UTC(),
	})
	return nil
}

func (service *UserService) UpdateUserDepartment(ctx context.Context, caller types.CallerContext, userID int, departmentID int) error {
	if !caller.IsAdmin() {
		service.logger.Warn("unauthorized attempt to update user department",
			slog.Int("caller_id", caller.UserID),
			slog.Int("user_id", userID),
		)
		return errors.ErrForbidden{Msg: "Only admins can move users between departments."}
	}

	user, err := service.repo.GetUserByID(ctx, userID)
	if err != nil {
		service.logger.Error("user not found when trying to update department",
			slog.Int("user_id", userID),
			slog.Int("caller_id", caller.UserID),
			slog.Any("error", err),
		)
		return err
	}

	if user.DepartmentID == nil {
		service.logger.Warn("cannot change department for non-department member",
			slog.Int("caller_id", caller.UserID),
			slog.Int("user_id", userID),
		)
		return errors.ErrBadRequest{Msg: "Cannot change department for a user who is not a department member."}
	}

	reqs, err := service.requestRepo.GetRequestsByDepartment(ctx, departmentID, nil)
	if err != nil {
		service.logger.Error("error when retrieving requests by department",
			slog.Int("caller_id", caller.UserID),
			slog.Any("error", err),
		)
		return err
	}

	for _, req := range reqs {
		if req.ClaimedBy != nil && *req.ClaimedBy == userID {
			service.logger.Warn("cannot move user who still has claimed requests",
				slog.Int("caller_id", caller.UserID),
				slog.Int("user_id", userID),
			)
			return errors.ErrBadRequest{Msg: "User must unclaim all requests before being moved to another department."}
		}
	}

	if err := service.repo.UpdateUserDepartment(ctx, userID, departmentID); err != nil {
		service.logger.Error("failed to update user department",
			slog.Int("user_id", userID),
			slog.Int("department_id", departmentID),
			slog.Int("caller_id", caller.UserID),
			slog.Any("error", err),
		)
		return err
	}

	service.logger.Info("user department updated successfully",
		slog.Int("user_id", userID),
		slog.Int("department_id", departmentID),
		slog.Int("caller_id", caller.UserID),
	)
	return nil
}

func (service *UserService) UpdateUserProfile(ctx context.Context, caller types.CallerContext, dto models.UserProfilePatchDTO) error {
	if err := service.repo.UpdateUserProfile(ctx, caller.UserID, dto); err != nil {
		service.logger.Error("failed to update user profile",
			slog.Int("caller_id", caller.UserID),
			slog.Any("error", err),
		)
		return errors.ErrInternalServerError{Msg: "Failed to update profile."}
	}

	service.logger.Info("user profile updated successfully",
		slog.Int("caller_id", caller.UserID),
	)
	return nil
}

func (service *UserService) DeactivateUser(ctx context.Context, caller types.CallerContext, id int) error {
	_, err := service.repo.GetUserByID(ctx, id)
	if err != nil {
		service.logger.Error("could not retrieve user from db",
			slog.Int("user_id", id),
			slog.Int("caller_id", caller.UserID),
			slog.Any("error", err),
		)
		return err
	}

	if !caller.IsAdmin() && !caller.IsDepartmentMember() {
		service.logger.Warn("unauthorized attempt to deactivate account",
			slog.Int("user_id", id),
			slog.Int("caller_id", caller.UserID),
		)
		return errors.ErrForbidden{Msg: "You do not have permission to deactivate this account."}
	}

	if err := service.repo.DeactivateUser(ctx, id); err != nil {
		service.logger.Error("error when trying to deactivate user",
			slog.Int("user_id", id),
			slog.Int("caller_id", caller.UserID),
			slog.Any("error", err),
		)
		return errors.ErrInternalServerError{Msg: fmt.Sprintf("Error deactivating user: %v", err)}
	}

	service.logger.Info("user deactivated successfully",
		slog.Int("user_id", id),
		slog.Int("caller_id", caller.UserID),
	)

	service.bus.Publish(ctx, events.Event{
		Type:         events.EventUserDeactivated,
		ActorID:      caller.UserID,
		ResourceID:   id,
		ResourceType: events.ResourceTypeUser,
		OccurredAt:   time.Now().UTC(),
	})
	return nil
}

func (service *UserService) GetUserByCognitoSub(ctx context.Context, cognitoSub string) (*models.User, error) {
	user, err := service.repo.GetUserByCognitoSub(ctx, cognitoSub)
	if err != nil {
		service.logger.Error("failed to fetch user by cognito sub",
			slog.String("cognito_sub", cognitoSub),
			slog.Any("error", err),
		)
		return nil, err
	}
	return &user, nil
}
