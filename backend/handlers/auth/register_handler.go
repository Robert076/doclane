package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/types/requests"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	req := requests.RegisterRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: fmt.Sprintf("Invalid body received. %v", err)})
		return
	}

	log.Print("Salut 1")
	id, err := config.UserService.AddUser(r.Context(), req.Email, req.Password, req.Role)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{Success: true, Data: id})
}
