package handlers

import (
	"context"
	"net/http"

	"github.com/xFidle/sportradar-intern/server/internal/httpx"
	"github.com/xFidle/sportradar-intern/server/internal/models"
)

type SportService interface {
	GetSports(ctx context.Context) ([]models.Sport, error)
}

type SportHandler struct {
	svc SportService
}

func NewSportHandler(svc SportService) *SportHandler {
	return &SportHandler{svc: svc}
}

func (h *SportHandler) HandleGetSports(w http.ResponseWriter, r *http.Request) {
	sports, err := h.svc.GetSports(context.TODO())
	if err != nil {
		logError(err, r)
		httpx.WriteError(w, httpx.InternalFailureError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, sports)
}
