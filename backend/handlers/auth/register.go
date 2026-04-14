// handlers/auth/register_handler.go
package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Robert076/doclane/backend/services"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

type RegisterRequest struct {
	Email          string  `json:"email"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	Password       string  `json:"password"`
	InvitationCode *string `json:"invitation_code"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	req := RegisterRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid JSON body format."})
		return
	}

	params := services.RegisterParams{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
		Role:      types.RoleMember,
	}

	var inviteCodeID *int

	if req.InvitationCode != nil && *req.InvitationCode != "" {
		invCode, err := config.InvitationCodeService.GetInvitationCodeInfo(r.Context(), *req.InvitationCode)
		if err != nil {
			utils.WriteError(w, err)
			return
		}
		params.Role = types.RoleMember
		params.DepartmentID = &invCode.DepartmentID
		inviteCodeID = &invCode.ID
	}

	id, err := config.UserService.AddUser(r.Context(), params)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	if inviteCodeID != nil {
		if err := config.InvitationCodeService.InvalidateCode(r.Context(), *inviteCodeID); err != nil {
			// user is created but code isn't invalidated — log it but don't fail the request
			// this is an edge case that can be handled by expiry
			_ = err
		}
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Msg:     "Registered successfully.",
		Data:    id,
	})
}
