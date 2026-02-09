package document_handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func PatchDocumentRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("\n\n1\n\n")
	userId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		log.Print("\n\n2\n\n")
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print("\n\n3\n\n")
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid document request ID format."})
		return
	}

	var dto models.DocumentRequestDTOPatch
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		log.Print("\n\n4\n\n")
		utils.WriteError(w, err)
		return
	}

	if err := config.DocumentService.PatchDocumentRequest(r.Context(), userId, id, dto); err != nil {
		log.Print("\n\n5\n\n")
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Document request updated successfully.",
		Data:    nil,
	})
}
