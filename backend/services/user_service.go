package services

import (
	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
)

type UserService struct {
	Repo repositories.IUserRepository
}

func (service *UserService) GetUserByEmail(email string) (models.User, error) {
	return service.Repo.GetUserByEmail(email)
}
