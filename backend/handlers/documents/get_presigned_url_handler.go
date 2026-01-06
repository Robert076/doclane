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

func GetFilePresignedURLHandler(w http.ResponseWriter, r *http.Request) {
	jwtUserId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	fileIDStr := chi.URLParam(r, "fileId")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid file ID format."})
		return
	}

	url, err := config.DocumentService.GetFilePresignedURL(r.Context(), jwtUserId, fileID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Presigned URL generated successfully.",
		Data:    map[string]string{"url": url},
	})
}
