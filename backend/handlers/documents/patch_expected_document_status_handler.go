package document_handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/requests"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func PatchExpectedDocumentStatusHandler(w http.ResponseWriter, r *http.Request) {
	docIDStr := chi.URLParam(r, "id")
	docID, err := strconv.Atoi(docIDStr)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	var req requests.UpdateExpectedDocumentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, err)
		return
	}

	if err = config.ExpectedDocumentService.UpdateExpectedDocumentStatus(r.Context(), docID, req.Status, req.RejectionReason); err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Document status updated successfully",
	})
}
