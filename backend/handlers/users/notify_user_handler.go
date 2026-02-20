package user_handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func NotifyUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	jwtUserID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	err = config.UserService.NotifyUser(r.Context(), jwtUserID, idInt)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	u, err := config.UserService.GetUserByID(r.Context(), idInt)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     fmt.Sprintf("%s %s was successfully notified.", u.FirstName, u.LastName),
	})
}
