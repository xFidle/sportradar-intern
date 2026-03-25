package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/xFidle/sportradar-intern/server/internal/models"
	"github.com/xFidle/sportradar-intern/server/internal/repo"
)

type CompetitionService struct {
	fAddr string
	db    *pgxpool.Pool
	q     *repo.Queries
}

func NewCompetitionService(db *pgxpool.Pool, fileserverAddr string) *CompetitionService {
	return &CompetitionService{
		fAddr: fileserverAddr,
		db:    db,
		q:     repo.New(db),
	}
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

	for i := range competitions {
		competitions[i].LogoPath = fmt.Sprintf("%s/%s", s.fAddr, competitions[i].LogoPath)
	}

	return competitions, nil
}
