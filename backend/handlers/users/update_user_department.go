package user_handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

type UpdateUserDepartmentRequest struct {
	DepartmentID int `json:"department_id"`
}

func UpdateUserDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.GetClaimsFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid user ID."})
		return
	}

	var req UpdateUserDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid request body."})
		return
	}

	if req.DepartmentID == 0 {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "department_id is required."})
		return
	}

	if err := config.UserService.UpdateUserDepartment(r.Context(), *claims, userID, req.DepartmentID); err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "User department updated successfully.",
	})
}
