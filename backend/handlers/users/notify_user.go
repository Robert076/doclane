package user_handler

import (
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

	userClaims, err := utils.GetCallerFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	err = config.UserService.NotifyUser(r.Context(), userClaims, idInt)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Notification was sent successfully.",
	})
}
