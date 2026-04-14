package invitation_handler

import (
	"net/http"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func GetInvitationCodeInfoHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Code is required."})
		return
	}

	invCode, err := config.InvitationCodeService.GetInvitationCodeInfo(r.Context(), code)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Invitation code is valid.",
		Data:    invCode,
	})
}
