package auth_middleware

import (
	"context"
	"net/http"

	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/golang-jwt/jwt"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := r.Cookie("auth_cookie")
		if err != nil {
			utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
			return
		}

		claims := &utils.CustomClaims{}
		token, err := jwt.ParseWithClaims(jwtToken.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
			return
		}

		ctx := context.WithValue(r.Context(), "jwtClaims", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
