package httpx

import (
	"log"
	"net/http"
)

var (
	InvalidPayloadError  = HError{Code: http.StatusBadRequest, Message: "invalid payload format"}
	InternalFailureError = HError{Code: http.StatusInternalServerError, Message: "internal processing failure"}
	ValidationError      = HError{Code: http.StatusUnprocessableEntity, Message: "invalid payload, validation failed"}
)

type HError struct {
	Code    int
	Message string
	Details any
}

func (h *HError) Error() string {
	return h.Message
}

func LogError(err error, r *http.Request) {
	log.Printf("%v: %v", r.URL.Path, err)
}
