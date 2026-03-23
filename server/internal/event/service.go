package event

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/xFidle/sportradar-intern/server/internal/models"
	"github.com/xFidle/sportradar-intern/server/internal/repo"
)

type Service struct {
	bucketAddr string
	db         *pgxpool.Pool
	q          *repo.Queries
}

func New(db *pgxpool.Pool) *Service {
	return &Service{
		db: db,
		q:  repo.New(db),
	}
}

func (s *Service) GetEvent(ctx context.Context, id int32) (*models.DetailedEvent, error) {
	return nil, nil
}

func (s *Service) GetEvents(ctx context.Context, filter models.Filter) ([]models.Event, error) {
	events, err := s.listBaseEvents(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := s.populateTeams(ctx, events); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Service) listBaseEvents(ctx context.Context, filter models.Filter) ([]models.Event, error) {
	params := repo.ListFilteredEventsParams{}
	if err := copier.Copy(&params, &filter); err != nil {
		return nil, err
	}

	layout := "2006-01-02"
	after, err := time.Parse(layout, filter.StartAfter)
	if err != nil {
		return nil, err
	}

	before, err := time.Parse(layout, filter.EndBefore)
	if err != nil {
		return nil, err
	}

	params.StartAfter = after
	params.EndBefore = before

	rows, err := s.q.ListFilteredEvents(ctx, params)
	if err != nil {
		return nil, err
	}

	var events []models.Event
	if err := copier.Copy(&events, &rows); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Service) populateTeams(ctx context.Context, events []models.Event) error {
	lookup := make(map[int32]*models.Event)
	for i := range events {
		lookup[events[i].EventID] = &events[i]
	}

	rows, err := s.q.ListEventTeams(ctx, slices.Collect(maps.Keys(lookup)))
	if err != nil {
		return err
	}

	for _, row := range rows {
		var team models.Team
		if err := copier.Copy(&team, &row); err != nil {
			return err
		}

		if e := lookup[row.EventID]; e != nil {
			e.Participants = append(e.Participants, team)
			fmt.Println(e.Participants)
		}
	}

	return nil
}
