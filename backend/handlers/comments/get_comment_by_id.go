package comment_handler

import (
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func GetCommentByID(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.GetClaimsFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	id := chi.URLParam(r, "commentID")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteError(w, errors.ErrUnprocessableContent{Msg: "Invalid ID received."})
		return
	}

	comm, err := config.RequestCommentService.GetCommentByID(r.Context(), *claims, idInt)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Comment retrieved successfully",
		Data:    comm,
	})
}
