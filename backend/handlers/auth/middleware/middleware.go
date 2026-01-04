package auth_middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/golang-jwt/jwt"
)

func AuthGuard(next http.Handler) http.Handler {
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

		uid, err := strconv.Atoi(claims.UserID)
		if err != nil {
			config.Logger.Error("failed to parse user id from claims",
				slog.String("user_id_raw", claims.UserID),
				slog.Any("error", err),
			)
			utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
			return
		}

		user, err := config.UserService.GetUserByID(r.Context(), uid)
		if err != nil {
			config.Logger.Error("database error in MustBeActive check",
				slog.Int("user_id", uid),
				slog.Any("error", err),
			)
			utils.WriteError(w, errors.ErrForbidden{Msg: "Access denied."})
			return
		}

		if !user.IsActive {
			config.Logger.Warn("deactivated user attempt to access protected route",
				slog.Int("user_id", uid),
				slog.String("email", user.Email),
			)
			utils.WriteError(w, errors.ErrForbidden{Msg: "Your account is deactivated."})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func MustBeProfessional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(utils.ClaimsKey).(*utils.CustomClaims)
		if !ok {
			config.Logger.Error("middleware context error: claims not found in MustBeProfessional")
			utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
			return
		}

		if claims.Role != "PROFESSIONAL" {
			config.Logger.Warn("unauthorized role access attempt",
				slog.String("user_id", claims.UserID),
				slog.String("required_role", "PROFESSIONAL"),
				slog.String("actual_role", claims.Role),
				slog.String("path", r.URL.Path),
			)
			utils.WriteError(w, errors.ErrForbidden{Msg: "Access denied. Professional role required."})
			return
		}

		next.ServeHTTP(w, r)
	})
}
