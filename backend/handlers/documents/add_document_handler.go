package document_handler

import (
	"net/http"
	"strconv"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/utils"
	"github.com/Robert076/doclane/backend/utils/config"
)

func AddDocumentHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		utils.WriteError(w, err)
		return
	}

	requestIDStr := r.FormValue("document_request_id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	defer file.Close()

	id, err := config.DocumentService.AddDocumentFile(
		r.Context(),
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
