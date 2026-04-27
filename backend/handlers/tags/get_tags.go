package tag_handler

import (
	"net/http"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func GetTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := config.TagService.GetTags(r.Context())
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
