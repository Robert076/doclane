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

func DeleteExpectedDocumentTemplateHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.GetClaimsFromContext(r.Context())
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

	expectedDocRequestTemplateIDStr := chi.URLParam(r, "expectedDocId")
	expectedDocRequestTemplateID, err := strconv.Atoi(expectedDocRequestTemplateIDStr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid expected document template ID format."})
		return
	}

	if err := config.RequestTemplateService.DeleteExpectedDocumentTemplate(
		r.Context(),
		*claims,
		expectedDocRequestTemplateID,
		templateID,
	); err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Expected document template deleted successfully.",
	})
}
