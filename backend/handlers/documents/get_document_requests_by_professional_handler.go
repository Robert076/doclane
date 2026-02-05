package document_handler

import (
	"net/http"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func GetDocumentRequestsByProfessionalHandler(w http.ResponseWriter, r *http.Request) {
	jwtUserId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	q := r.URL.Query()
	var searchPtr *string

	// search
	if s := q.Get("search"); s != "" {
		searchPtr = &s
	}

	reqs, err := config.DocumentService.GetDocumentRequestsByProfessional(r.Context(), jwtUserId, searchPtr)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Professional document requests retrieved successfully.",
		Data:    reqs,
	})
}
