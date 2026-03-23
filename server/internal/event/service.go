package event

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xFidle/sportradar-intern/server/internal/models"
	"github.com/xFidle/sportradar-intern/server/internal/util"
)

const layout = "2006-01-02"

type Service struct {
	loader *loader
}

func New(db *pgxpool.Pool, fileserverAddr string) *Service {
	return &Service{loader: newLoader(db, fileserverAddr)}
}

func (s *Service) GetEvent(ctx context.Context, id int32) (*models.DetailedEvent, error) {
	event, err := s.loader.fetchEventByID(ctx, id)
	if err != nil {
		return nil, err
	}

	teams, err := s.loader.fetchTeamsByEventID(ctx, id)
	if err != nil {
		return nil, err
	}

	teamIDs := util.Map(teams, func(t models.DetailedTeam) int32 { return t.TeamID })
	players, err := s.loader.fetchPlayersByTeamID(ctx, teamIDs)
	if err != nil {
		return nil, err
	}

	lookup := make(map[int32]*models.DetailedTeam)
	for i := range teams {
		lookup[teams[i].TeamID] = &teams[i]
	}

	for _, p := range players {
		if t := lookup[p.teamID]; t != nil {
			t.Players = append(t.Players, p.player)
		}
	}

	event.Participants = teams

	return event, nil
}

func (s *Service) GetEvents(ctx context.Context, filter models.Filter) ([]models.Event, error) {
	events, err := s.loader.fetchEventsByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	eventIDs := util.Map(events, func(e models.Event) int32 { return e.EventID })
	teams, err := s.loader.fetchTeamsByEventIDs(ctx, eventIDs)
	if err != nil {
		return nil, err
	}

	lookup := make(map[int32]*models.Event)
	for i := range events {
		lookup[events[i].EventID] = &events[i]
	}

	for _, t := range teams {
		if e := lookup[t.eventID]; e != nil {
			e.Participants = append(e.Participants, t.team)
		}
	}

	return events, nil
}
