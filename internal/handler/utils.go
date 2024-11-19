package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
)

func isValidID(w http.ResponseWriter, id string, place, ops string) bool {
	ValidBucketPath := regexp.MustCompile(`^order[1-9][0-9]*$`)
	if !ValidBucketPath.MatchString(id) {
		s := "Invalid ID in: " + place
		slog.Error(s)
		s = "Failed to " + ops + " order"
		writeError(w, s, http.StatusNotFound)
		return false
	}

	return true
}

func isJSONFile(w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		slog.Error("Invalid content type: expected application/json")
		writeError(w, "Content type must be 'application/json'", http.StatusBadRequest)
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
	// slog.Warn(message)
}
