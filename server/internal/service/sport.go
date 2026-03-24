package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/xFidle/sportradar-intern/server/internal/models"
	"github.com/xFidle/sportradar-intern/server/internal/repo"
)

type SportService struct {
	db *pgxpool.Pool
	q  *repo.Queries
}

func NewSportService(db *pgxpool.Pool) *SportService {
	return &SportService{db: db, q: repo.New(db)}
}

func (s *SportService) GetSports(ctx context.Context) ([]models.Sport, error) {
	rows, err := s.q.ListSports(ctx)
	if err != nil {
		return nil, err
	}

	var sports []models.Sport
	if err := copier.Copy(&sports, &rows); err != nil {
		return nil, err
	}

	return sports, nil
}
