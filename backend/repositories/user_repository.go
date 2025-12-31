package repositories

import "github.com/Robert076/doclane/backend/models"

type UserRepository struct {
}

func (repo *UserRepository) GetUserByEmail(email string) (models.User, error) {
	return models.User{}, nil
}
