package event

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xFidle/sportradar-intern/server/internal/models"
)

type Service struct {
	db *pgxpool.Pool
}

func (s *Service) GetEvent(id int) (*models.DetailedEvent, error) {
	return nil, nil
}

func (s *Service) GetEvents(filter models.Filter) ([]models.Event, error) {
	return nil, nil
}
