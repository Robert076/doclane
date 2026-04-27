package services

import (
	"strings"

	apierrors "github.com/Robert076/doclane/backend/types/errors"
)

func validateTagDTO(name, color string) error {
	if strings.TrimSpace(name) == "" {
		return apierrors.ErrBadRequest{Msg: "Tag name is required."}
	}
	if color != "" && !isValidHexColor(color) {
		return apierrors.ErrBadRequest{Msg: "Tag color must be a valid hex color (e.g. #ff5722)."}
	}
	return nil
}

func isValidHexColor(s string) bool {
	if len(s) != 7 || s[0] != '#' {
		return false
	}
	for _, c := range s[1:] {
		if !('0' <= c && c <= '9') && !('a' <= c && c <= 'f') && !('A' <= c && c <= 'F') {
			return false
		}
	}
	return true
}
