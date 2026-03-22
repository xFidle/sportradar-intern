package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/xFidle/sportradar-intern/server/internal/httpx"
	"github.com/xFidle/sportradar-intern/server/internal/models"
)

type EventService interface {
	GetEvent(ctx context.Context, id int32) (*models.DetailedEvent, error)
	GetEvents(ctx context.Context, filter models.Filter) ([]models.Event, error)
}

type EventHandler struct {
	svc EventService
}

func NewEventHandler(svc EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

func (h *EventHandler) HandleGetEvent(w http.ResponseWriter, r *http.Request) {

}

func (h *EventHandler) HandleGetEvents(w http.ResponseWriter, r *http.Request) {
	var filter models.Filter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		httpx.LogError(err, r)
		httpx.WriteError(w, httpx.InvalidPayloadError)
		return
	}

	v := httpx.NewValdiator()
	if err := v.Struct(filter); err != nil {
		httpx.LogError(err, r)
		httpx.WriteError(w, httpx.ExtractValidationError(err))
		return
	}

	events, err := h.svc.GetEvents(context.TODO(), filter)
	if err != nil {
		httpx.LogError(err, r)
		httpx.WriteError(w, httpx.InternalFailureError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, events)
}
