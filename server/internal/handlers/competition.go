package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/xFidle/sportradar-intern/server/internal/httpx"
	"github.com/xFidle/sportradar-intern/server/internal/models"
)

type CompetitionService interface {
	GetCompetitionsBySportID(ctx context.Context, id int32) ([]models.Competition, error)
}

type CompetitionHandler struct {
	svc CompetitionService
}

func NewCompetitionHandler(svc CompetitionService) *CompetitionHandler {
	return &CompetitionHandler{svc: svc}
}

func (h *CompetitionHandler) HandleGetCompetitions(w http.ResponseWriter, r *http.Request) {
	sportID, err := strconv.Atoi(r.URL.Query().Get("sport_id"))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InvalidPathParameter)
		return
	}

	competitions, err := h.svc.GetCompetitionsBySportID(context.TODO(), int32(sportID))
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InternalFailureError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, competitions)
}
