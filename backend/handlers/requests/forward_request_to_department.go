package request_handler

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

func ForwardRequestToDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.GetClaimsFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	requestIDStr := chi.URLParam(r, "id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid request ID format."})
		return
	}

	var body struct {
		DepartmentID int `json:"department_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid request body."})
		return
	}

	if body.DepartmentID == 0 {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "department_id is required."})
		return
	}

	if err := config.RequestService.ForwardRequestToDepartment(r.Context(), *claims, requestID, body.DepartmentID); err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Request forwarded to department successfully.",
	})
}
