package utils

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/utils/awscfg"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/golang-jwt/jwt"
)

type contextKey string

const ClaimsKey contextKey = "jwtClaims"

var JWTSecret string

func init() {
	var jwt string
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		awsCfg := awscfg.InitAWSConfig()
		ssmClient := ssm.NewFromConfig(awsCfg)

		jwtPath := os.Getenv("JWT_SECRET_PATH")
		jwtParam, err := ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{
			Name:           aws.String(jwtPath),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			log.Fatalf("error when trying to get jwt from SSM %v", err)
		}
		jwt = *jwtParam.Parameter.Value
	} else {
		jwt = RequireEnv("JWT_SECRET")
	}
	JWTSecret = jwt
}

func GenerateJWT(user models.User) (string, error) {
	claims := types.JWTClaims{
		UserID:       user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
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
