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

func GetDocumentRequestsByClientHandler(w http.ResponseWriter, r *http.Request) {
	jwtUserId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	clientIDStr := chi.URLParam(r, "clientID")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid client ID format."})
		return
	}

	reqs, err := config.DocumentService.GetDocumentRequestsByClient(r.Context(), jwtUserId, clientID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Client document requests retrieved successfully.",
		Data:    reqs,
	})
}
