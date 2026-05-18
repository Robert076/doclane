package user_handler

import (
	"encoding/json"
	"net/http"

	"github.com/Robert076/doclane/backend/services"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

type syncUserRequest struct {
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	InvitationCode *string `json:"invitation_code"`
}

func SyncUserHandler(w http.ResponseWriter, r *http.Request) {
	caller, err := utils.GetCallerFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	var req syncUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid JSON body."})
		return
	}

	// Determine role — default to member unless the seed secret is provided.
	// This allows seeding an admin account for demo purposes without a
	// dedicated handler. The secret is an env var and never exposed publicly.
	role := types.RoleMember
	seedSecret := utils.RequireEnv("SEED_SECRET")
	if r.Header.Get("X-Seed-Secret") == seedSecret && r.Header.Get("X-Role") == types.RoleAdmin {
		role = types.RoleAdmin
	}

	params := services.SyncUserParams{
		CognitoSub: caller.CognitoSub,
		Email:      r.Context().Value(types.CognitoEmailKey).(string),
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

		_ = config.InvitationCodeService.InvalidateCode(r.Context(), invCode.ID)

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
