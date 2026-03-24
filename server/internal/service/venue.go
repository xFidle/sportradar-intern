package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/xFidle/sportradar-intern/server/internal/models"
	"github.com/xFidle/sportradar-intern/server/internal/repo"
)

type VenueService struct {
	db *pgxpool.Pool
	q  *repo.Queries
}

func NewVenueService(db *pgxpool.Pool) *VenueService {
	return &VenueService{db: db, q: repo.New(db)}
}

func (s *VenueService) GetVenuesBySportID(ctx context.Context, id int32) ([]models.Venue, error) {
	rows, err := s.q.ListVenuesBySportID(ctx, id)
	if err != nil {
		return nil, err
	}

	var venues []models.Venue
	if err := copier.Copy(&venues, &rows); err != nil {
		return nil, err
	}

	for i := range venues {
		venues[i].City = models.City{
			Name:    rows[i].CityName,
			Country: models.Country{Name: rows[i].CountryName, Code: rows[i].CountryCode},
		}
	}

	return venues, err
}
