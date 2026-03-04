package template_handler

import (
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func PresignExampleHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	templateIDStr := chi.URLParam(r, "id")
	templateID, err := strconv.Atoi(templateIDStr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid template ID format."})
		return
	}

	expectedDocTemplateIDStr := chi.URLParam(r, "expectedDocId")
	expectedDocTemplateID, err := strconv.Atoi(expectedDocTemplateIDStr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid expected document template ID format."})
		return
	}

	url, err := config.DocumentRequestTemplateService.PresignExample(r.Context(), userID, templateID, expectedDocTemplateID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Data: url,
	})
}
