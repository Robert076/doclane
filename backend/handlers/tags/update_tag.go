package tag_handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func UpdateTagHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.GetClaimsFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid tag ID."})
		return
	}

	var dto models.TagDTOUpdate
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid request body."})
		return
	}

	tag, err := config.TagService.UpdateTag(r.Context(), *claims, id, dto)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Tag updated successfully.",
		Data:    tag,
	})
}
