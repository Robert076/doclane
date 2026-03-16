package template_handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func InstantiateRequestTemplateHandler(w http.ResponseWriter, r *http.Request) {
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

	var body struct {
		ClientID     int        `json:"client_id"`
		IsScheduled  bool       `json:"is_scheduled"`
		ScheduledFor *string    `json:"scheduled_for"`
		DueDate      *time.Time `json:"due_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid request body."})
		return
	}

	if body.ClientID == 0 {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "client_id is required."})
		return
	}

	if body.IsScheduled && body.ScheduledFor == nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "scheduled_for is required when is_scheduled is true."})
		return
	}

	id, err := config.RequestTemplateService.InstantiateRequestTemplate(
		r.Context(),
		userID,
		templateID,
		body.ClientID,
		body.IsScheduled,
		body.ScheduledFor,
		body.DueDate,
	)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Msg:     "RequestTemplate instantiated successfully.",
		Data:    id,
	})
}
