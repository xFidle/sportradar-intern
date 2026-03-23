package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/xFidle/sportradar-intern/server/internal/httpx"
	"github.com/xFidle/sportradar-intern/server/internal/models"
	"github.com/xFidle/sportradar-intern/server/internal/util"
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
	eventID, err := strconv.Atoi(r.URL.Query().Get("event_id"))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InvalidPayloadError)
		return
	}

	event, err := h.svc.GetEvent(context.TODO(), int32(eventID))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InternalFailureError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, event)
}

func (h *EventHandler) HandleGetEvents(w http.ResponseWriter, r *http.Request) {
	var filter models.Filter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InvalidPayloadError)
		return
	}

	v := httpx.NewValdiator()
	if err := v.Struct(filter); err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.ExtractValidationError(err))
		return
	}

	if _, err := time.Parse(util.DateLayout, filter.StartAfter); err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InvalidPayloadError)
		return
	}

	if _, err := time.Parse(util.DateLayout, filter.EndBefore); err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InvalidPayloadError)
		return
	}

	events, err := h.svc.GetEvents(context.TODO(), filter)
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InternalFailureError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, events)
}
