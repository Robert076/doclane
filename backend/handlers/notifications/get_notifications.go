package notification_handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.GetCallerFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	notifications, err := config.AuditLogService.GetNotifications(r.Context(), claims, limit)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	user, _ := config.UserService.GetUserByCognitoSub(r.Context(), claims.CognitoSub)
	var seenAt *time.Time
	if user != nil {
		seenAt = user.NotificationsSeenAt
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Notifications retrieved successfully.",
		Data: map[string]any{
			"notifications": notifications,
			"seen_at":       seenAt,
		},
	})
}
