package handlers

import (
	"net/http"

	"github.com/xFidle/sportradar-intern/server/internal/event"
	"github.com/xFidle/sportradar-intern/server/internal/models"
)

type EventService interface {
	GetEvent(id int) (*models.DetailedEvent, error)
	GetEvents(filter models.Filter) ([]models.Event, error)
}

type EventHandler struct {
	svc event.Service
}

func (h *EventHandler) HandleGetEvent(w http.ResponseWriter, r *http.Request) {

}

func (h *EventHandler) HandleGetEvents(w http.ResponseWriter, r *http.Request) {

}
