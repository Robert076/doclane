package utils

import (
	"net/http"

	"github.com/Robert076/doclane/backend/types"
	"github.com/Robert076/doclane/backend/types/errors"
)

func WriteError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case errors.ErrNotFound:
		WriteJSONSafe(w, http.StatusNotFound, types.APIResponse{Success: false, Err: err.Error()})
	case errors.ErrBadRequest:
		WriteJSONSafe(w, http.StatusBadRequest, types.APIResponse{Success: false, Err: err.Error()})
	case errors.ErrFileTypeNotSupported:
		WriteJSONSafe(w, http.StatusUnsupportedMediaType, types.APIResponse{Success: false, Err: err.Error()})
	case errors.ErrFileSizeTooBig:
		WriteJSONSafe(w, http.StatusRequestEntityTooLarge, types.APIResponse{Success: false, Err: err.Error()})
	case errors.ErrUnauthorized:
		WriteJSONSafe(w, http.StatusUnauthorized, types.APIResponse{Success: false, Err: err.Error()})
	case errors.ErrForbidden:
		WriteJSONSafe(w, http.StatusForbidden, types.APIResponse{Success: false, Err: err.Error()})
	case errors.ErrConflict:
		WriteJSONSafe(w, http.StatusConflict, types.APIResponse{Success: false, Err: err.Error()})
	case errors.ErrUnprocessableContent:
		WriteJSONSafe(w, http.StatusUnprocessableEntity, types.APIResponse{Success: false, Err: err.Error()})
	case errors.ErrTooManyRequests:
		WriteJSONSafe(w, http.StatusTooManyRequests, types.APIResponse{Success: false, Err: err.Error()})
	case errors.ErrInternalServerError:
		WriteJSONSafe(w, http.StatusInternalServerError, types.APIResponse{Success: false, Err: err.Error()})
	case errors.ErrBadGateway:
		WriteJSONSafe(w, http.StatusBadGateway, types.APIResponse{Success: false, Err: err.Error()})
	default:
		WriteJSONSafe(w, http.StatusBadRequest, types.APIResponse{Success: false, Err: err.Error()})
	}
}
