package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Robert076/doclane/backend/services"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/types/requests"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	req := requests.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid request body."})
		return
	}

	params := services.LoginParams{
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := config.UserService.Login(r.Context(), params)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		utils.WriteError(w, errors.ErrInternalServerError{Msg: "Could not generate session."})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_cookie",
		Value:    token,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Expires:  time.Now().Add(time.Hour * 24),
	})

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Login successful.",
		Data:    token,
	})
}
