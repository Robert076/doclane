package invitation_handler

import (
	"encoding/json"
	"net/http"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func ValidateInvitationCodeHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.InvitationCodeValidateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid request body."})
		return
	}

	if dto.Code == "" {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Code is required."})
		return
	}

	invCode, err := config.InvitationCodeService.GetInvitationCodeByCode(
		r.Context(),
		dto.Code,
	)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Invitation code is valid.",
		Data:    map[string]int{"professional_id": invCode.ProfessionalID},
	})
}
