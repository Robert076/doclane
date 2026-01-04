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

func GetFilesByRequestHandler(w http.ResponseWriter, r *http.Request) {
	jwtUserId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	requestIDStr := chi.URLParam(r, "requestID")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid request ID format."})
		return
	}

	files, err := config.DocumentService.GetFilesByRequest(r.Context(), jwtUserId, requestID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Files retrieved successfully.",
		Data:    files,
	})
}
