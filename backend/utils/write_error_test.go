package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrors "github.com/Robert076/doclane/backend/types/errors"
)

func decodeBody(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	return body
}

func TestWriteError_NotFoundMapsTo404(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, apierrors.ErrNotFound{Msg: "user not found"})

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	body := decodeBody(t, rec)
	if body["success"] != false {
		t.Errorf("expected success=false, got %v", body["success"])
	}
	if body["error"] != "user not found" {
		t.Errorf("expected error message to be passed through, got %v", body["error"])
	}
}

func TestWriteError_BadRequestMapsTo400(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, apierrors.ErrBadRequest{Msg: "missing field"})

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestWriteError_UnauthorizedMapsTo401(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, apierrors.ErrUnauthorized{Msg: "Unauthorized."})

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestWriteError_ForbiddenMapsTo403(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, apierrors.ErrForbidden{Msg: "not allowed"})

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}

func TestWriteError_ConflictMapsTo409(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, apierrors.ErrConflict{Msg: "already exists"})

	if rec.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
}

func TestWriteError_FileTypeNotSupportedMapsTo415(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, apierrors.ErrFileTypeNotSupported{Msg: "only pdf allowed"})

	if rec.Code != http.StatusUnsupportedMediaType {
		t.Errorf("expected status %d, got %d", http.StatusUnsupportedMediaType, rec.Code)
	}
}

func TestWriteError_FileSizeTooBigMapsTo413(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, apierrors.ErrFileSizeTooBig{Msg: "file too large"})

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected status %d, got %d", http.StatusRequestEntityTooLarge, rec.Code)
	}
}

func TestWriteError_UnprocessableContentMapsTo422(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, apierrors.ErrUnprocessableContent{Msg: "cannot process"})

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestWriteError_TooManyRequestsMapsTo429(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, apierrors.ErrTooManyRequests{Msg: "slow down"})

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("expected status %d, got %d", http.StatusTooManyRequests, rec.Code)
	}
}

func TestWriteError_InternalServerErrorMapsTo500(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, apierrors.ErrInternalServerError{Msg: "boom"})

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestWriteError_BadGatewayMapsTo502(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, apierrors.ErrBadGateway{Msg: "upstream failed"})

	if rec.Code != http.StatusBadGateway {
		t.Errorf("expected status %d, got %d", http.StatusBadGateway, rec.Code)
	}
}

func TestWriteError_UnknownErrorDefaultsTo400(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, errors.New("some unexpected error"))

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected default status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	body := decodeBody(t, rec)
	if body["error"] != "some unexpected error" {
		t.Errorf("expected error message to be passed through, got %v", body["error"])
	}
}
