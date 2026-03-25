package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/xFidle/sportradar-intern/server/internal/httpx"
	"github.com/xFidle/sportradar-intern/server/internal/models"
)

type TeamService interface {
	GetTeamsBySportID(ctx context.Context, id int32) ([]models.Team, error)
	GetTeamsByCompetitionID(ctx context.Context, id int32) ([]models.Team, error)
}

type TeamHandler struct {
	svc TeamService
}

func NewTeamHandler(svc TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) HandleGetTeamsBySport(w http.ResponseWriter, r *http.Request) {
	sport_id, err := strconv.Atoi(chi.URLParam(r, "sport_id"))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InvalidPathParameter)
		return
	}

	teams, err := h.svc.GetTeamsBySportID(context.TODO(), int32(sport_id))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InternalFailureError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, teams)
}

func (h *TeamHandler) HandleGetTeamsByCompetition(w http.ResponseWriter, r *http.Request) {
	sport_id, err := strconv.Atoi(chi.URLParam(r, "competition_id"))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InvalidPathParameter)
		return
	}

	teams, err := h.svc.GetTeamsByCompetitionID(context.TODO(), int32(sport_id))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InternalFailureError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, teams)
}
