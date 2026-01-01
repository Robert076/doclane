package services

import (
	"github.com/Robert076/doclane/backend/models"
)

type IUserService interface {
	AddUser(email string, password string) (int, error)
	GetUserByEmail(email string) (models.User, error)
}
