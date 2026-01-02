package repositories

import "github.com/Robert076/doclane/backend/models"

type IUserRepository interface {
	GetUsers(limit *int, offset *int, orderBy *string, order *string) ([]models.User, error)
	AddUser(user models.User) (int, error)
	GetUserByEmail(email string) (models.User, error)
}
