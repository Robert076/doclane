package document_handler

import (
	"encoding/json"
	"net/http"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/types/requests"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func AddDocumentRequestHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	req := requests.AddDocumentRequestRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, err)
		return
	}

	id, err := config.DocumentService.AddDocumentRequest(r.Context(), userId, req.ProfessionalID, req.ClientID, req.Title, req.Description, req.DueDate)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Document request created successfully.",
		Data:    id,
	})
}
