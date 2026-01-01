package repositories

import "github.com/Robert076/doclane/backend/models"

type IUserRepository interface {
	AddUser(user models.User) (int, error)
	GetUserByEmail(email string) (models.User, error)
}
