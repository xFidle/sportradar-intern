package sport

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/xFidle/sportradar-intern/server/internal/models"
	"github.com/xFidle/sportradar-intern/server/internal/repo"
)

type Service struct {
	db *pgxpool.Pool
	q  *repo.Queries
}

func New(db *pgxpool.Pool) *Service {
	return &Service{db: db, q: repo.New(db)}
}

func (s *Service) GetSports(ctx context.Context) ([]models.Sport, error) {
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
