package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/xFidle/sportradar-intern/server/internal/models"
	"github.com/xFidle/sportradar-intern/server/internal/repo"
)

type CompetitionService struct {
	db *pgxpool.Pool
	q  *repo.Queries
}

func NewCompetitionService(db *pgxpool.Pool) *CompetitionService {
	return &CompetitionService{db: db, q: repo.New(db)}
}

func (s *CompetitionService) GetCompetitionsBySportID(ctx context.Context, id int32) ([]models.Competition, error) {
	rows, err := s.q.ListCompetitionsBySportID(ctx, id)
	if err != nil {
		return nil, err
	}

	var competitions []models.Competition
	if err := copier.Copy(&competitions, &rows); err != nil {
		return nil, err
	}

	return competitions, nil
}
