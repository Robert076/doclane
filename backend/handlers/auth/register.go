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
	DepartmentID   *int    `json:"department_id"`
	InvitationCode *string `json:"invitation_code"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	req := RegisterRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid JSON body format."})
		return
	}

	params := services.RegisterParams{
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Password:     req.Password,
		Role:         types.RoleMember,
		DepartmentID: req.DepartmentID,
	}

	id, err := config.UserService.AddUser(r.Context(), params)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Msg:     "Registered successfully.",
		Data:    id,
	})
}
