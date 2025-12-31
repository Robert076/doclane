package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	req := types.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := config.UserService.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusBadRequest)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "invalid email or password", http.StatusBadRequest)
		return
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		http.Error(w, "could not generate jwt token", http.StatusInternalServerError)
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

	w.WriteHeader(http.StatusOK)
}
