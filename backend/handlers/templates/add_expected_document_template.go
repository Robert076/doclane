package template_handler

import (
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

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
		RequestTemplateID: templateID,
		Title:             title,
		Description:       description,
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

	id, err := config.RequestTemplateService.AddExpectedDocumentTemplate(
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
