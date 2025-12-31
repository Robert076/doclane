package config

import (
	"log"
	"os"

	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/services"
)

var JWTSecret string
var userRepository *repositories.UserRepository
var UserService *services.UserService

func init() {
	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	userRepository = &repositories.UserRepository{}
	UserService = &services.UserService{Repo: userRepository}
}
