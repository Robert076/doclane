package request_handler

import (
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/events"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func GetRequestAuditLogHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.GetCallerFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid request ID format."})
		return
	}

	if _, err := config.RequestService.GetRequestByID(r.Context(), claims, id); err != nil {
		utils.WriteError(w, err)
		return
	}

	auditEvents, err := config.AuditLogService.GetByResource(r.Context(), events.ResourceTypeRequest, id)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Audit log retrieved successfully.",
		Data:    auditEvents,
	})
}
