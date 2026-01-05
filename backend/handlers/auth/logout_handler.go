package auth

import (
	"net/http"
	"time"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/utils"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_cookie",
		Value:    "",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Expires:  time.Now(),
	})

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Logout successful.",
	})
}
