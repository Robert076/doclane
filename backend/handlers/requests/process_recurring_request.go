package request_handler

import (
	"net/http"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func ProcessRecurringRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if err := config.RequestService.ProcessRecurringRequests(r.Context()); err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Recurring requests processed successfully.",
	})
}
