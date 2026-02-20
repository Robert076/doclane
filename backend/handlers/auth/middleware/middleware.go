package auth_middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/golang-jwt/jwt"
)

func AuthGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// first check Bearer token, then cookie if that token is not present.
		var tokenString string

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			}
		}

		if tokenString == "" {
			jwtCookie, err := r.Cookie("auth_cookie")
			if err != nil {
				utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
				return
			}
			tokenString = jwtCookie.Value
		}

		claims := &utils.CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			config.Logger.Error("invalid jwt token",
				slog.Any("error", err),
				slog.String("remote_addr", r.RemoteAddr),
			)
			utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
			return
		}

		ctx := context.WithValue(r.Context(), utils.ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MustBeActive(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(utils.ClaimsKey).(*utils.CustomClaims)
		if !ok {
			config.Logger.Error("middleware context error: claims not found")
			utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
			return
		}

		user, err := config.UserService.GetUserByID(r.Context(), claims.UserID)
		if err != nil {
			config.Logger.Error("database error in MustBeActive check",
				slog.Int("user_id", claims.UserID),
				slog.Any("error", err),
			)
			utils.WriteError(w, errors.ErrForbidden{Msg: "Access denied."})
			return
		}

		if !user.IsActive {
			config.Logger.Warn("deactivated user attempt to access protected route",
				slog.Int("user_id", claims.UserID),
				slog.String("email", user.Email),
			)
			utils.WriteError(w, errors.ErrForbidden{Msg: "Your account is deactivated."})
			return
		}

		next.ServeHTTP(w, r)
	})
}
