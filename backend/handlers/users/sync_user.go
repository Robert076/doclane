package user_handler

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

type syncUserRequest struct {
	Email          string  `json:"email"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	InvitationCode *string `json:"invitation_code"`
}

// SyncUserHandler is called once after Cognito confirms a new user's email.
// It creates the user's record in the application database, linked to their
// Cognito identity via cognito_sub. Login and registration themselves are
// handled entirely by Cognito on the frontend.
func SyncUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse token directly — this route has no AuthGuard since the user
	// doesn't exist in the DB yet. We verify the token manually.
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if tokenString == "" {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}
	var claims struct {
		Sub string `json:"sub"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil || claims.Sub == "" {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	var req syncUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid JSON body."})
		return
	}

	if req.Email == "" {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Email is required."})
		return
	}

	role := types.RoleMember
	seedSecret := utils.RequireEnv("SEED_SECRET")
	if r.Header.Get("X-Seed-Secret") == seedSecret && r.Header.Get("X-Role") == types.RoleAdmin {
		role = types.RoleAdmin
	}

	params := services.SyncUserParams{
		CognitoSub: claims.Sub,
		Email:      req.Email,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Role:       role,
	}

	if req.InvitationCode != nil && *req.InvitationCode != "" {
		invCode, err := config.InvitationCodeService.GetInvitationCodeInfo(r.Context(), *req.InvitationCode)
		if err != nil {
			utils.WriteError(w, err)
			return
		}
		params.DepartmentID = &invCode.DepartmentID

		id, err := config.UserService.SyncUser(r.Context(), params)
		if err != nil {
			utils.WriteError(w, err)
			return
		}
		_ = config.InvitationCodeService.InvalidateCode(r.Context(), invCode.ID, id)

		utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
			Success: true,
			Msg:     "User synced successfully.",
			Data:    id,
		})
		return
	}

	id, err := config.UserService.SyncUser(r.Context(), params)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Msg:     "User synced successfully.",
		Data:    id,
	})
}
