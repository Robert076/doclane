package template_handler

import (
	"fmt"
	"net/http"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func AddRequestTemplateWithDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	const maxRequestSize = 100 << 20 // 100MB to allow multiple files
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)

	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Failed to parse form."})
		return
	}
	defer r.MultipartForm.RemoveAll()

	title := r.FormValue("title")
	description := r.FormValue("description")
	isRecurring := r.FormValue("is_recurring") == "true"
	recurrenceCron := r.FormValue("recurrence_cron")

	if title == "" {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Title is required."})
		return
	}

	template := models.RequestTemplate{
		Title:       title,
		Description: &description,
		IsRecurring: isRecurring,
	}
	if recurrenceCron != "" {
		template.RecurrenceCron = &recurrenceCron
	}

	// Parse expected documents from form fields:
	// expected_documents[0][title], expected_documents[0][description]
	// expected_documents[0][example_file] (optional file)
	var docs []types.ExpectedDocumentTemplateInput
	for i := 0; ; i++ {
		docTitle := r.FormValue(fmt.Sprintf("expected_documents[%d][title]", i))
		if docTitle == "" {
			break
		}
		docDescription := r.FormValue(fmt.Sprintf("expected_documents[%d][description]", i))

		input := types.ExpectedDocumentTemplateInput{
			Title:       docTitle,
			Description: docDescription,
		}

		file, header, err := r.FormFile(fmt.Sprintf("expected_documents[%d][example_file]", i))
		if err == nil {
			defer file.Close()
			input.ExampleFile = file
			input.ExampleFileName = header.Filename
			input.ExampleMimeType = header.Header.Get("Content-Type")
			input.ExampleFileSize = header.Size
		}

		docs = append(docs, input)
	}

	id, err := config.RequestTemplateService.AddRequestTemplateWithDocuments(r.Context(), userID, template, docs)
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
