package utils

import (
	"errors"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	UserID         string `json:"user_id"`
	Role           string `json:"role"`
	ProfessionalID string `json:"professional_id,omitempty"`
	jwt.StandardClaims
}

func GenerateJWT(user models.User) (string, error) {
	profID := ""
	if user.ProfessionalID != nil {
		profID = *user.ProfessionalID
	}

	claims := CustomClaims{
		UserID:         user.ID,
		Role:           user.Role,
		ProfessionalID: profID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "doclane.app",
			Audience:  "doclane.app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.JWTSecret))
}

func ValidateJWT(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token provided.")
	}

	return claims, nil
}
