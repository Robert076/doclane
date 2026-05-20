package notification_handler

import (
	"net/http"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func MarkNotificationsSeenHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.GetCallerFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	if err := config.UserService.MarkNotificationsSeen(r.Context(), claims); err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Notifications marked as seen.",
	})
}
