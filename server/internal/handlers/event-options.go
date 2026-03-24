package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/xFidle/sportradar-intern/server/internal/httpx"
	"github.com/xFidle/sportradar-intern/server/internal/models"
)

type CompetitionService interface {
	GetCompetitionsBySportID(ctx context.Context, id int32) ([]models.Competition, error)
}

type VenueService interface {
	GetVenuesBySportID(ctx context.Context, id int32) ([]models.Venue, error)
}

type EventOptionsHandler struct {
	cSvc CompetitionService
	vSvc VenueService
}

func NewEventOptionsHandler(cSvc CompetitionService, vSvc VenueService) *EventOptionsHandler {
	return &EventOptionsHandler{
		cSvc: cSvc,
		vSvc: vSvc,
	}
}

func (h *EventOptionsHandler) HandleGetCompetitionsBySport(w http.ResponseWriter, r *http.Request) {
	sportID, err := strconv.Atoi(chi.URLParam(r, "sport_id"))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InvalidPathParameter)
		return
	}

	competitions, err := h.cSvc.GetCompetitionsBySportID(context.TODO(), int32(sportID))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InternalFailureError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, competitions)
}

type EventOptions struct {
	Venues       []models.Venue       `json:"venues"`
	Competitions []models.Competition `json:"competitions"`
}

func (h *EventOptionsHandler) HandleGetEventOptionsBySport(w http.ResponseWriter, r *http.Request) {
	sportID, err := strconv.Atoi(chi.URLParam(r, "sport_id"))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InvalidPathParameter)
		return
	}

	venues, err := h.vSvc.GetVenuesBySportID(context.TODO(), int32(sportID))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InternalFailureError)
		return
	}

	competitions, err := h.cSvc.GetCompetitionsBySportID(context.TODO(), int32(sportID))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InternalFailureError)
		return
	}

	options := EventOptions{Venues: venues, Competitions: competitions}

	httpx.WriteJSON(w, http.StatusOK, options)
}
