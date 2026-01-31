package document_handler

import (
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
	"github.com/go-chi/chi/v5"
)

func AddDocumentHandler(w http.ResponseWriter, r *http.Request) {
	const maxRequestSize = 21 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)

	userId, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, errors.ErrUnauthorized{Msg: "Unauthorized."})
		return
	}

	requestIDStr := chi.URLParam(r, "id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Invalid request ID format."})
		return
	}

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils.WriteError(w, err)
		return
	}
	defer func() {
		if err := r.MultipartForm.RemoveAll(); err != nil {
			utils.WriteError(w, errors.ErrInternalServerError{Msg: "Error removing temp files from multipart form"})
		}
	}()

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, errors.ErrBadRequest{Msg: "Could not get file from request."})
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			utils.WriteError(w, errors.ErrInternalServerError{Msg: "Error closing file"})
		}
	}()

	id, err := config.DocumentService.AddDocumentFile(
		r.Context(),
		userId,
		requestID,
		header.Filename,
		header.Size,
		header.Header.Get("Content-Type"),
		file,
	)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSONSafe(w, http.StatusCreated, types.APIResponse{
		Success: true,
		Msg:     "Successfully uploaded document.",
		Data:    id,
	})
}
