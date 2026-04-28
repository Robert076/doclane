package request_handler

import (
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func ExtractFileTextHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.GetClaimsFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	fileID, err := strconv.Atoi(chi.URLParam(r, "fileId"))
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid file ID."})
		return
	}

	text, err := config.RequestService.ExtractFileText(r.Context(), *claims, fileID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Text extracted successfully.",
		Data:    map[string]string{"text": text},
	})
}
