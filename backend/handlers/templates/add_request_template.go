package template_handler

import (
	"encoding/json"
	"net/http"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func AddRequestTemplateHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	var body struct {
		Title          string  `json:"title"`
		Description    *string `json:"description"`
		IsRecurring    bool    `json:"is_recurring"`
		RecurrenceCron *string `json:"recurrence_cron"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid request body."})
		return
	}

	if body.Title == "" {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Title is required."})
		return
	}

	template := models.RequestTemplate{
		Title:          body.Title,
		Description:    body.Description,
		IsRecurring:    body.IsRecurring,
		RecurrenceCron: body.RecurrenceCron,
	}

	id, err := config.RequestTemplateService.AddRequestTemplate(r.Context(), userID, template)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Msg:     "RequestTemplate created successfully.",
		Data:    id,
	})
}
