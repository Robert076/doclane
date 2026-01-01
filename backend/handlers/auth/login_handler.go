package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/types/requests"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	req := requests.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: fmt.Sprintf("Invalid request body received. %v", err)})
		return
	}

	user, err := config.UserService.GetUserByEmail(req.Email)
	if err != nil {
		utils.WriteError(w, errors.ErrNotFound{Msg: "Invalid email or password."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid email or password."})
		return
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		utils.WriteError(w, errors.ErrInternalServerError{Msg: fmt.Sprintf("Could not encode token. %v", err)})
		return
	}

	authCookie := http.Cookie{
		Name:     "auth_cookie",
		Value:    token,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Expires:  time.Now().Add(time.Hour * 24),
	}

	http.SetCookie(w, &authCookie)

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{Success: true, Token: token})
}
