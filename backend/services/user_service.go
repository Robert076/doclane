package services

import (
	"fmt"
	"net/mail"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repositories.IUserRepository
}

func NewUserService(repo repositories.IUserRepository) *UserService {
	return &UserService{repo: repo}
}

func (service *UserService) GetUsers(limit *int, offset *int, orderBy *string, order *string) ([]models.User, error) {
	return service.repo.GetUsers(limit, offset, orderBy, order)
}

func (service *UserService) AddUser(email string, password string, role string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	if err := service.ValidateUserForRegister(email, password); err != nil {
		return 0, err
	}

	user := models.User{Email: email, PasswordHash: string(hashedPassword), Role: role}
	return service.repo.AddUser(user)
}

func (service *UserService) GetUserByEmail(email string) (models.User, error) {
	return service.repo.GetUserByEmail(email)
}

func (service *UserService) ValidateUserForRegister(email string, password string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.ErrBadRequest{Msg: fmt.Sprintf("Invalid email received: %s", email)}
	}
	_, err := service.repo.GetUserByEmail(email)
	if err == nil {
		return errors.ErrConflict{Msg: fmt.Sprintf("User already exists.")}
	}
	if !errors.IsNotFound(err) {
		return errors.ErrInternalServerError{Msg: fmt.Sprintf("Failed to check if user already exists. %v", err)}
	}

	return nil
}
