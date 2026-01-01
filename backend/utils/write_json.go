package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}
	return nil
}

func WriteJSONSafe(w http.ResponseWriter, status int, data any) {
	if err := WriteJSON(w, status, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
