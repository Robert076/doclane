package invitation_handler

import (
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func GetInvitationCodesByDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.GetClaimsFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	deptIDStr := r.URL.Query().Get("department_id")
	if deptIDStr == "" {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Missing department_id query parameter."})
		return
	}

	deptID, err := strconv.Atoi(deptIDStr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid department_id value."})
		return
	}

	codes, err := config.InvitationCodeService.GetInvitationCodesByDepartment(r.Context(), *claims, deptID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Invitation codes retrieved successfully.",
		Data:    codes,
	})
}
