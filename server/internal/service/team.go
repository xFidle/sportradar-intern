package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/xFidle/sportradar-intern/server/internal/models"
	"github.com/xFidle/sportradar-intern/server/internal/repo"
)

type TeamService struct {
	fAddr string
	db    *pgxpool.Pool
	q     *repo.Queries
}

func NewTeamService(db *pgxpool.Pool, fileserverAddr string) *TeamService {
	return &TeamService{
		fAddr: fileserverAddr,
		db:    db,
		q:     repo.New(db),
	}
}

func (s *TeamService) GetTeamsByCompetitionID(ctx context.Context, id int32) ([]models.Team, error) {
	rows, err := s.q.ListTeamsBySportID(ctx, id)
	if err != nil {
		return nil, err
	}

	var teams []models.Team
	if err := copier.Copy(&teams, &rows); err != nil {
		return nil, err
	}

	for i := range teams {
		teams[i].LogoPath = fmt.Sprintf("%s/%s", s.fAddr, teams[i].LogoPath)
	}

	return teams, nil
}

func (s *TeamService) GetTeamsBySportID(ctx context.Context, id int32) ([]models.Team, error) {
	rows, err := s.q.ListTeamsBySportID(ctx, id)
	if err != nil {
		return nil, err
	}

	var teams []models.Team
	if err := copier.Copy(&teams, &rows); err != nil {
		return nil, err
	}

	for i := range teams {
		teams[i].LogoPath = fmt.Sprintf("%s/%s", s.fAddr, teams[i].LogoPath)
	}

	return teams, nil
}
