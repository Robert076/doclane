package user_handler

import (
	"net/http"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	user, err := config.UserService.GetUserByID(r.Context(), userId)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Data:    user,
		Msg:     "Retrieved logged in user successfully.",
	})
}
