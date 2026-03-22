package handlers

import "net/http"

var (
	InvalidPayloadError  = HTTPError{Code: http.StatusBadRequest, Message: "invalid payload format"}
	InternalFailureError = HTTPError{Code: http.StatusInternalServerError, Message: "internal processing failure"}
)

type HTTPError struct {
	Code    int
	Message string
}

func (h *HTTPError) Error() string {
	return h.Message
}
