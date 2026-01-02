package services

import (
	"github.com/Robert076/doclane/backend/models"
)

type IUserService interface {
	GetUsers(limit *int, offset *int, orderBy *string, order *string) ([]models.User, error)
	GetUserByEmail(email string) (models.User, error)
	AddUser(email string, password string) (int, error)
}
