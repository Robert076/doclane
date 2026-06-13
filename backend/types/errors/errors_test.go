package errors

import (
	"errors"
	"testing"
)

func TestIsNotFound_TrueForErrNotFound(t *testing.T) {
	if !IsNotFound(ErrNotFound{Msg: "missing"}) {
		t.Error("expected IsNotFound to return true for ErrNotFound")
	}
}

func TestIsNotFound_FalseForOtherError(t *testing.T) {
	if IsNotFound(ErrBadRequest{Msg: "bad"}) {
		t.Error("expected IsNotFound to return false for ErrBadRequest")
	}
}

func TestIsNotFound_FalseForPlainError(t *testing.T) {
	if IsNotFound(errors.New("plain")) {
		t.Error("expected IsNotFound to return false for a plain error")
	}
}

func TestIsBadRequest_TrueForErrBadRequest(t *testing.T) {
	if !IsBadRequest(ErrBadRequest{Msg: "bad"}) {
		t.Error("expected IsBadRequest to return true for ErrBadRequest")
	}
}

func TestIsBadRequest_FalseForOtherError(t *testing.T) {
	if IsBadRequest(ErrNotFound{Msg: "missing"}) {
		t.Error("expected IsBadRequest to return false for ErrNotFound")
	}
}

func TestIsUnauthorized_TrueForErrUnauthorized(t *testing.T) {
	if !IsUnauthorized(ErrUnauthorized{Msg: "nope"}) {
		t.Error("expected IsUnauthorized to return true for ErrUnauthorized")
	}
}

func TestIsForbidden_TrueForErrForbidden(t *testing.T) {
	if !IsForbidden(ErrForbidden{Msg: "denied"}) {
		t.Error("expected IsForbidden to return true for ErrForbidden")
	}
}

func TestIsConflict_TrueForErrConflict(t *testing.T) {
	if !IsConflict(ErrConflict{Msg: "exists"}) {
		t.Error("expected IsConflict to return true for ErrConflict")
	}
}

func TestIsUnprocessableContent_TrueForErrUnprocessableContent(t *testing.T) {
	if !IsUnprocessableContent(ErrUnprocessableContent{Msg: "nope"}) {
		t.Error("expected IsUnprocessableContent to return true for ErrUnprocessableContent")
	}
}

func TestIsTooManyRequests_TrueForErrTooManyRequests(t *testing.T) {
	if !IsTooManyRequests(ErrTooManyRequests{Msg: "slow"}) {
		t.Error("expected IsTooManyRequests to return true for ErrTooManyRequests")
	}
}

func TestIsInternalServerError_TrueForErrInternalServerError(t *testing.T) {
	if !IsInternalServerError(ErrInternalServerError{Msg: "boom"}) {
		t.Error("expected IsInternalServerError to return true for ErrInternalServerError")
	}
}

func TestIsBadGateway_TrueForErrBadGateway(t *testing.T) {
	if !IsBadGateway(ErrBadGateway{Msg: "upstream"}) {
		t.Error("expected IsBadGateway to return true for ErrBadGateway")
	}
}

func TestError_ReturnsMessage(t *testing.T) {
	err := ErrNotFound{Msg: "user 42 not found"}
	if err.Error() != "user 42 not found" {
		t.Errorf("expected Error() to return the message, got %q", err.Error())
	}
}
