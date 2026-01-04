package document_handler

import (
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func GetDocumentRequestByIDHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid document request ID format."})
		return
	}

	docReq, err := config.DocumentService.GetDocumentRequestByID(r.Context(), userId, id)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Document request retrieved successfully.",
		Data:    docReq,
	})
}
