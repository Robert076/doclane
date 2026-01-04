package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Robert076/doclane/backend/services"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/types/requests"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	req := requests.RegisterRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid JSON body format."})
		return
	}

	params := services.CreateUserParams{
		Email:          req.Email,
		Password:       req.Password,
		Role:           req.Role,
		ProfessionalID: req.ProfessionalID,
	}

	id, err := config.UserService.AddUser(r.Context(), params)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Data:    id,
	})
}
