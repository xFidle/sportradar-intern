package httpx

import (
	"net/http"
)

var (
	InvalidPathParameter = HError{
		Code:    http.StatusBadRequest,
		Message: "invalid path parameter",
	}
	InvalidPayloadError = HError{
		Code:    http.StatusBadRequest,
		Message: "invalid payload format",
	}
	InternalFailureError = HError{
		Code:    http.StatusInternalServerError,
		Message: "internal processing failure",
	}
	ValidationError = HError{
		Code:    http.StatusUnprocessableEntity,
		Message: "invalid payload, validation failed",
	}
)

type HError struct {
	Code    int
	Message string
	Details any
}

func (h *HError) Error() string {
	return h.Message
}
