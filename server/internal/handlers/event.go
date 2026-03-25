package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/xFidle/sportradar-intern/server/internal/httpx"
	"github.com/xFidle/sportradar-intern/server/internal/models"
	"github.com/xFidle/sportradar-intern/server/internal/service"
	"github.com/xFidle/sportradar-intern/server/internal/util"
)

type EventService interface {
	CreateEvent(ctx context.Context, req models.CreateEventReq) (*models.DetailedEvent, error)
	GetEvent(ctx context.Context, id int32) (*models.DetailedEvent, error)
	GetEvents(ctx context.Context, filter models.Filter) ([]models.Event, error)
}

type EventHandler struct {
	svc EventService
}

func NewEventHandler(svc EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

func (h *EventHandler) HandlePostEvent(w http.ResponseWriter, r *http.Request) {
	var createEventReq models.CreateEventReq
	if err := json.NewDecoder(r.Body).Decode(&createEventReq); err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InvalidPayloadError)
		return
	}

	v := httpx.NewValdiator()
	if err := v.Struct(createEventReq); err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.ExtractValidationError(err))
		return
	}

	t, err := time.Parse(util.TimestampLayout, createEventReq.StartTime)
	if err != nil {
		logError(err, r)
		hErr := httpx.InvalidPayloadError
		hErr.Details = fmt.Sprintf(util.ParsingTimestampMessage, createEventReq.StartTime)
		httpx.WriteError(w, hErr)
		return
	}

	if t.Before(time.Now()) {
		hErr := httpx.InvalidPayloadError
		hErr.Details = fmt.Sprintf("timestamp %s cannot be in past", createEventReq.StartTime)
		httpx.WriteError(w, hErr)
		return
	}

	if len(createEventReq.TeamIDs) < 2 || !util.AreUnique(createEventReq.TeamIDs) {
		hErr := httpx.InvalidPayloadError
		hErr.Details = "event needs at least 2 unique participants"
		httpx.WriteError(w, hErr)
		return
	}

	event, err := h.svc.CreateEvent(context.TODO(), createEventReq)
	if err != nil {
		logError(err, r)
		var hErr httpx.HError
		if errors.Is(err, service.ErrInvalidVenue) || errors.Is(err, service.ErrInvalidTeams) {
			hErr = httpx.InvalidPayloadError
			hErr.Details = err.Error()
		} else {
			hErr = httpx.InternalFailureError
		}
		httpx.WriteError(w, hErr)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, event)
}

func (h *EventHandler) HandleGetEvent(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.Atoi(chi.URLParam(r, "event_id"))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InvalidPathParameter)
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

	start, err := time.Parse(util.DateLayout, filter.StartAfter)
	if err != nil {
		logError(err, r)
		hErr := httpx.InvalidPayloadError
		hErr.Details = fmt.Sprintf(util.ParsingDateMessage, filter.StartAfter)
		httpx.WriteError(w, hErr)
		return
	}

	end, err := time.Parse(util.DateLayout, filter.EndBefore)
	if err != nil {
		logError(err, r)
		hErr := httpx.InvalidPayloadError
		hErr.Details = fmt.Sprintf(util.ParsingDateMessage, filter.EndBefore)
		httpx.WriteError(w, hErr)
		return
	}

	if start.After(end) {
		hErr := httpx.InvalidPayloadError
		hErr.Details = "'start_after' must be lower than 'end_before'"
		httpx.WriteError(w, hErr)
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
