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

func GetDocumentRequestsByProfessionalHandler(w http.ResponseWriter, r *http.Request) {
	jwtUserId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	profIDStr := chi.URLParam(r, "professionalID")
	profID, err := strconv.Atoi(profIDStr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid professional ID format."})
		return
	}

	reqs, err := config.DocumentService.GetDocumentRequestsByProfessional(r.Context(), jwtUserId, profID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Professional document requests retrieved successfully.",
		Data:    reqs,
	})
}
