package insertadmin_handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Robert076/doclane/backend/services"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

// InsertAdminHandler seeds a demo admin account directly into the DB.
// It expects a valid Cognito token in the request — the admin must already
// exist in Cognito before calling this. Used for demo/development only.
func InsertAdminHandler(w http.ResponseWriter, r *http.Request) {
	// Verify seed secret first
	seedSecret := utils.RequireEnv("SEED_SECRET")
	if r.Header.Get("X-Seed-Secret") != seedSecret {
		utils.WriteError(w, errors.ErrForbidden{Msg: "Forbidden."})
		return
	}

	// Extract sub and email directly from the Authorization header
	// without doing a DB lookup (user doesn't exist yet)
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if tokenString == "" {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	// Parse claims without verification — seed secret already proves intent,
	// and this endpoint is only called once manually
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid token."})
		return
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid token."})
		return
	}
	var claims struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil || claims.Sub == "" {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid token claims."})
		return
	}

	params := services.SyncUserParams{
		CognitoSub: claims.Sub,
		Email:      claims.Email,
		FirstName:  "Adminescu",
		LastName:   "Adminovici",
		Role:       types.RoleAdmin,
	}

	id, err := config.UserService.SyncUser(r.Context(), params)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Msg:     "Admin seeded successfully.",
		Data:    id,
	})
}
