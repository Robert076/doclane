// handlers/auth/register_handler.go
package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Robert076/doclane/backend/services"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/types/requests"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func RegisterProfessionalHandler(w http.ResponseWriter, r *http.Request) {
	req := requests.RegisterProfessionalRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid JSON body format."})
		return
	}

	params := services.CreateUserParams{
		Email:          req.Email,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Password:       req.Password,
		Role:           "PROFESSIONAL",
		ProfessionalID: nil,
	}

	id, err := config.UserService.AddUser(r.Context(), params)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Msg:     "Professional registered successfully.",
		Data:    id,
	})
}

func RegisterClientHandler(w http.ResponseWriter, r *http.Request) {
	req := requests.RegisterClientRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid JSON body format."})
		return
	}

	if req.InvitationCode == "" {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invitation code is required."})
		return
	}

	profID, err := config.InvitationCodeService.ValidateAndUseInvitationCode(r.Context(), req.InvitationCode)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	params := services.CreateUserParams{
		Email:          req.Email,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Password:       req.Password,
		Role:           "CLIENT",
		ProfessionalID: &profID,
	}

	id, err := config.UserService.AddUser(r.Context(), params)
	if err != nil {
		if reactivateErr := config.InvitationCodeService.ReactivateCode(r.Context(), req.InvitationCode); reactivateErr != nil {
			utils.WriteError(w, errors.ErrBadRequest{Msg: fmt.Sprintf("%v , %v", err, reactivateErr)})
			return
		}
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Msg:     "Client registered successfully.",
		Data:    id,
	})
}
