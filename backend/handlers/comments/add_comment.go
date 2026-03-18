package comment_handler

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

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	reqID := chi.URLParam(r, "id")
	reqIDInt, err := strconv.Atoi(reqID)
	if err != nil {
		utils.WriteError(w, errors.ErrUnprocessableContent{Msg: "Invalid ID received."})
		return
	}

	var comm models.RequestComment
	if err := json.NewDecoder(r.Body).Decode(&comm); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "The body received is invalid."})
		return
	}

	id, err := config.RequestCommentService.AddComment(r.Context(), userID, reqIDInt, comm)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Comment added successfully.",
		Data:    id,
	})
}
