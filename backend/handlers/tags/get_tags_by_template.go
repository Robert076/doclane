package tag_handler

import (
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func GetTagsByTemplateIDHandler(w http.ResponseWriter, r *http.Request) {
	templateID, err := strconv.Atoi(chi.URLParam(r, "templateId"))
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid template ID."})
		return
	}

	tags, err := config.TagService.GetTagsByTemplateID(r.Context(), templateID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Tags retrieved successfully.",
		Data:    tags,
	})
}
