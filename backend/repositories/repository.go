package repositories

import "github.com/Robert076/doclane/backend/models"

type IUserRepository interface {
	GetUserByEmail(email string) (models.User, error)
}
