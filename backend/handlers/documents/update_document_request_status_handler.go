package document_handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/types/requests"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func UpdateDocumentRequestStatusHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIDFromContext(r.Context())
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

	var req requests.UpdateDocumentRequestStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid JSON body."})
		return
	}

	err = config.DocumentService.UpdateDocumentRequestStatus(r.Context(), userId, id, req.Status)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Document request status updated successfully.",
	})
}
