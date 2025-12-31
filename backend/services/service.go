package services

import (
	"github.com/Robert076/doclane/backend/models"
)

type IUserService interface {
	GetUserByEmail(email string) (models.User, error)
}
