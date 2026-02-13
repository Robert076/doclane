package user_handler

import (
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func DeactivateUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "You must provide an user id for deactivation."})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid user id received for deactivation."})
		return
	}

	if err := config.UserService.DeactivateUser(r.Context(), userId, idInt); err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Client successfully deactivated",
	})
}
