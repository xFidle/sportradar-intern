package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/xFidle/sportradar-intern/server/internal/httpx"
	"github.com/xFidle/sportradar-intern/server/internal/models"
)

type TeamService interface {
	GetTeamsByCompetitionID(ctx context.Context, id int32) ([]models.Team, error)
}

type TeamHandler struct {
	svc TeamService
}

func NewTeamHandler(svc TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) HandleGetTeamsByCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := strconv.Atoi(r.URL.Query().Get("competition_id"))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InvalidPathParameter)
		return
	}

	teams, err := h.svc.GetTeamsByCompetitionID(context.TODO(), int32(competitionID))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InternalFailureError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, teams)
}
