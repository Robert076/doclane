package document_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func AddRequestHandler(w http.ResponseWriter, r *http.Request) {
	const maxRequestSize = 21 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)

	userId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	contentType := r.Header.Get("Content-Type")

	var dto models.RequestDTOCreate

	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(5 << 20); err != nil {
			utils.WriteError(w, errors.ErrBadRequest{Msg: "Failed to parse form."})
			return
		}
		defer r.MultipartForm.RemoveAll()

		dto.Title = r.FormValue("title")
		dto.ClientID, err = strconv.Atoi(r.FormValue("client_id"))
		if err != nil {
			utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid client_id."})
			return
		}

		description := r.FormValue("description")
		if description != "" {
			dto.Description = &description
		}

		if r.FormValue("is_recurring") == "true" {
			dto.IsRecurring = true
		}

		recurrenceCron := r.FormValue("recurrence_cron")
		if recurrenceCron != "" {
			dto.RecurrenceCron = &recurrenceCron
		}

		if r.FormValue("is_scheduled") == "true" {
			dto.IsScheduled = true
		}

		scheduledFor := r.FormValue("scheduled_for")
		if scheduledFor != "" {
			dto.ScheduledFor = &scheduledFor
		}

		dueDateStr := r.FormValue("due_date")
		if dueDateStr != "" {
			dueDate, err := time.Parse(time.RFC3339, dueDateStr)
			if err != nil {
				utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid due_date format."})
				return
			}
			dto.DueDate = &dueDate
		}

		i := 0
		for {
			title := r.FormValue(fmt.Sprintf("expected_documents[%d][title]", i))
			if title == "" {
				break
			}
			description := r.FormValue(fmt.Sprintf("expected_documents[%d][description]", i))

			ed := models.ExpectedDocumentInput{
				Title:       title,
				Description: description,
			}

			file, header, err := r.FormFile(fmt.Sprintf("expected_documents[%d][example_file]", i))
			if err == nil {
				ed.ExampleFile = file
				ed.ExampleFileName = header.Filename
				ed.ExampleMimeType = header.Header.Get("Content-Type")
				ed.ExampleFileSize = header.Size
			}

			dto.ExpectedDocuments = append(dto.ExpectedDocuments, ed)
			i++
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			utils.WriteError(w, err)
			return
		}
	}

	id, err := config.RequestService.AddRequest(r.Context(), userId, dto)
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
