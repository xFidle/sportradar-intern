package httpx

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

type ResponseError struct {
	Error   string `json:"error"`
	Details any    `json:"details"`
}

func WriteError(w http.ResponseWriter, err HError) {
	WriteJSON(w, err.Code, ResponseError{Error: err.Error(), Details: err.Details})
}
