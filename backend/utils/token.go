package utils

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type contextKey string

const ClaimsKey contextKey = "jwtClaims"

var JWTSecret string

func init() {
	godotenv.Load("../../.env")
	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET env var is not set")
	}
}

func GenerateJWT(user models.User) (string, error) {
	claims := types.JWTClaims{
		UserID:       user.ID,
		Role:         user.Role,
		DepartmentID: user.DepartmentID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "doclane.app",
			Audience:  "doclane.app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(JWTSecret))
}

func ValidateJWT(tokenString string) (*types.JWTClaims, error) {
	claims := &types.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token provided.")
	}

	return claims, nil
}

func GetClaimsFromContext(ctx context.Context) (*types.JWTClaims, error) {
	claims, ok := ctx.Value(ClaimsKey).(*types.JWTClaims)
	if !ok {
		return nil, errors.New("could not find user claims in context")
	}
	return claims, nil
}

func GetUserIDFromContext(ctx context.Context) (int, error) {
	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}
