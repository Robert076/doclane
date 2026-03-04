package template_handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func AddTemplateHandler(w http.ResponseWriter, r *http.Request) {
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

	template := models.DocumentRequestTemplate{
		Title:          body.Title,
		Description:    body.Description,
		IsRecurring:    body.IsRecurring,
		RecurrenceCron: body.RecurrenceCron,
	}

	id, err := config.DocumentRequestTemplateService.AddTemplate(r.Context(), userID, template)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Msg:     "Template created successfully.",
		Data:    id,
	})
}

func GetTemplatesByProfessionalHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	templates, err := config.DocumentRequestTemplateService.GetTemplatesByProfessionalID(r.Context(), userID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Templates retrieved successfully.",
		Data:    templates,
	})
}

func GetTemplateByIDHandler(w http.ResponseWriter, r *http.Request) {
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

	template, err := config.DocumentRequestTemplateService.GetTemplateByID(r.Context(), userID, templateID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusOK, types.APIResponse{
		Success: true,
		Msg:     "Template retrieved successfully.",
		Data:    template,
	})
}

func AddExpectedDocumentTemplateHandler(w http.ResponseWriter, r *http.Request) {
	const maxRequestSize = 21 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)

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

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Failed to parse form."})
		return
	}
	defer r.MultipartForm.RemoveAll()

	title := r.FormValue("title")
	description := r.FormValue("description")

	if title == "" || description == "" {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Title and description are required."})
		return
	}

	t := models.ExpectedDocumentTemplate{
		DocumentRequestTemplateID: templateID,
		Title:                     title,
		Description:               description,
	}

	var (
		exampleFile     interface{ Read([]byte) (int, error) } = nil
		exampleFileName string
		ExampleMimeType string
		exampleFileSize int64
	)

	file, header, err := r.FormFile("example_file")
	if err == nil {
		defer file.Close()
		exampleFile = file
		exampleFileName = header.Filename
		ExampleMimeType = header.Header.Get("Content-Type")
		exampleFileSize = header.Size
	}

	id, err := config.DocumentRequestTemplateService.AddExpectedDocumentTemplate(
		r.Context(),
		userID,
		t,
		exampleFile,
		exampleFileName,
		ExampleMimeType,
		exampleFileSize,
	)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Msg:     "Expected document template added successfully.",
		Data:    id,
	})
}

func DeleteExpectedDocumentTemplateHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := config.DocumentRequestTemplateService.DeleteExpectedDocumentTemplate(
		r.Context(),
		userID,
		expectedDocTemplateID,
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

func InstantiateTemplateHandler(w http.ResponseWriter, r *http.Request) {
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

	id, err := config.DocumentRequestTemplateService.InstantiateTemplate(
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
		Msg:     "Template instantiated successfully.",
		Data:    id,
	})
}
