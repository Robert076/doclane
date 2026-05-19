package request_handler

import (
	"encoding/json"
	"net/http"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func AddRequestHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.GetCallerFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	var dto models.RequestDTOCreate
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid request body."})
		return
	}

	id, err := config.RequestService.AddRequest(r.Context(), claims, dto)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Msg:     "Document request created successfully.",
		Data:    id,
	})
}
