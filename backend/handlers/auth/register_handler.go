package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/types/requests"
	"github.com/Robert076/doclane/backend/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	req := requests.RegisterRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: fmt.Sprintf("Invalid body received. %v", err)})
		return
	}
}
