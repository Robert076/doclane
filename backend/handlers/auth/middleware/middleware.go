package auth_middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/golang-jwt/jwt/v5"
)

// jwksCache caches Cognito's public keys so we don't fetch them on every request.
// Cognito rotates keys rarely — a 1 hour cache is safe.
var (
	jwksCache     []jwk
	jwksCacheTime time.Time
	jwksMu        sync.RWMutex
)

// jwk represents a single JSON Web Key from Cognito's JWKS endpoint.
type jwk struct {
	Kid string `json:"kid"` // Key ID — matches the "kid" header in the token
	N   string `json:"n"`   // RSA modulus (base64url encoded)
	E   string `json:"e"`   // RSA exponent (base64url encoded)
}

// getPublicKey fetches Cognito's JWKS and finds the RSA public key
// matching the given kid (key ID from the token header).
func getPublicKey(region, userPoolID, kid string) (*rsa.PublicKey, error) {
	jwksMu.RLock()
	if jwksCache != nil && time.Since(jwksCacheTime) < time.Hour {
		key, err := findKey(jwksCache, kid)
		jwksMu.RUnlock()
		return key, err
	}
	jwksMu.RUnlock()

	url := fmt.Sprintf(
		"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
		region, userPoolID,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Keys []jwk `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	jwksMu.Lock()
	jwksCache = result.Keys
	jwksCacheTime = time.Now()
	jwksMu.Unlock()

	return findKey(result.Keys, kid)
}

// findKey finds the JWK matching the given kid and converts it to an *rsa.PublicKey.
func findKey(keys []jwk, kid string) (*rsa.PublicKey, error) {
	for _, k := range keys {
		if k.Kid != kid {
			continue
		}

		nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
		if err != nil {
			return nil, fmt.Errorf("failed to decode modulus: %w", err)
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
		if err != nil {
			return nil, fmt.Errorf("failed to decode exponent: %w", err)
		}

		n := new(big.Int).SetBytes(nBytes)
		e := new(big.Int).SetBytes(eBytes)

		return &rsa.PublicKey{
			N: n,
			E: int(e.Int64()),
		}, nil
	}
	return nil, fmt.Errorf("no key found for kid: %s", kid)
}

// AuthGuard verifies the Cognito-issued JWT and resolves the caller's identity
// from the application database. On success it stores a CallerContext in the
// request context under types.CallerContextKey.
//
// Passwords and credentials are never seen here — Cognito owns those.
func AuthGuard(region, userPoolID, clientID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// 1. Extract token from Authorization header or cookie
			var tokenString string
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			}
			if tokenString == "" {
				cookie, err := r.Cookie("auth_cookie")
				if err != nil {
					utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
					return
				}
				tokenString = cookie.Value
			}

			// 2. Decode token header (unverified) to get the kid
			parts := strings.Split(tokenString, ".")
			if len(parts) != 3 {
				utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
				return
			}
			headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
			if err != nil {
				utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
				return
			}
			var header struct {
				Kid string `json:"kid"`
			}
			if err := json.Unmarshal(headerJSON, &header); err != nil || header.Kid == "" {
				utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
				return
			}

			// 3. Fetch the matching RSA public key from Cognito (cached)
			pubKey, err := getPublicKey(region, userPoolID, header.Kid)
			if err != nil {
				config.Logger.Error("failed to get Cognito public key",
					slog.String("kid", header.Kid),
					slog.Any("error", err),
				)
				utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
				return
			}

			// 4. Verify the token signature
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return pubKey, nil
			})
			if err != nil || !token.Valid {
				config.Logger.Error("invalid Cognito token",
					slog.Any("error", err),
					slog.String("remote_addr", r.RemoteAddr),
				)
				utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
				return
			}

			// 5. Verify token was issued for this app
			if aud, _ := claims["aud"].(string); aud != clientID {
				utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
				return
			}

			// 6. Extract the Cognito sub
			sub, _ := claims["sub"].(string)
			if sub == "" {
				utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
				return
			}

			// 7. Resolve sub to a DB user
			user, err := config.UserService.GetUserByCognitoSub(r.Context(), sub)
			if err != nil {
				config.Logger.Error("failed to resolve Cognito sub to DB user",
					slog.String("cognito_sub", sub),
					slog.Any("error", err),
				)
				utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
				return
			}

			// 8. Store CallerContext in request context
			caller := types.CallerContext{
				UserID:       user.ID,
				CognitoSub:   sub,
				Role:         user.Role,
				DepartmentID: user.DepartmentID,
			}
			ctx := context.WithValue(r.Context(), types.CallerContextKey, caller)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// MustBeActive rejects requests from deactivated users.
// It must be chained after AuthGuard since it reads CallerContext.
func MustBeActive(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		caller, ok := r.Context().Value(types.CallerContextKey).(types.CallerContext)
		if !ok {
			config.Logger.Error("MustBeActive: CallerContext missing — AuthGuard must run first")
			utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
			return
		}

		user, err := config.UserService.GetUserByID(r.Context(), caller, caller.UserID)
		if err != nil {
			config.Logger.Error("database error in MustBeActive check",
				slog.Int("user_id", caller.UserID),
				slog.Any("error", err),
			)
			utils.WriteError(w, errors.ErrForbidden{Msg: "Access denied."})
			return
		}

		if !user.IsActive {
			config.Logger.Warn("deactivated user attempted to access protected route",
				slog.Int("user_id", caller.UserID),
				slog.String("email", user.Email),
			)
			utils.WriteError(w, errors.ErrForbidden{Msg: "Your account is deactivated."})
			return
		}

		next.ServeHTTP(w, r)
	})
}
