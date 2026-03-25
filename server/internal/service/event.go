package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/xFidle/sportradar-intern/server/internal/models"
	"github.com/xFidle/sportradar-intern/server/internal/repo"
	"github.com/xFidle/sportradar-intern/server/internal/util"
)

var (
	ErrInvalidVenue = errors.New("venue does not belong to competition")
	ErrInvalidTeams = errors.New("one or more teams do not belong to competition")
)

type EventService struct {
	loader *loader
	db     *pgxpool.Pool
	q      *repo.Queries
}

func NewEventService(db *pgxpool.Pool, fileserverAddr string) *EventService {
	return &EventService{
		loader: newLoader(db, fileserverAddr),
		db:     db,
		q:      repo.New(db),
	}
}

func (s *EventService) CreateEvent(ctx context.Context, req models.CreateEventReq) (*models.DetailedEvent, error) {
	isVenueCorrect, err := s.loader.checkVenueCorrectnes(ctx, req.CompetitionID, req.VenueID)
	if err != nil {
		return nil, err
	}

	if !isVenueCorrect {
		return nil, ErrInvalidVenue
	}

	areTeamsCorrect, err := s.loader.checkTeamsCorrectness(ctx, req.CompetitionID, req.TeamIDs)
	if err != nil {
		return nil, err
	}

	if !areTeamsCorrect {
		return nil, ErrInvalidTeams
	}

	return s.insertEventTX(ctx, req)
}

func (s *EventService) GetEvent(ctx context.Context, id int32) (*models.DetailedEvent, error) {
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

	if event.Status != models.Scheduled {
		scores, err := s.loader.fetchScoresByEventID(ctx, id)
		if err != nil {
			return nil, err
		}
		event.Scores = scores
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

func (s *EventService) GetEvents(ctx context.Context, filter models.Filter) ([]models.Event, error) {
	events, err := s.loader.fetchEventsByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	eventIDs := util.Map(events, func(e models.Event) int32 { return e.EventID })
	teams, err := s.loader.fetchTeamsByEventIDs(ctx, eventIDs)
	if err != nil {
		return nil, err
	}

	scores, err := s.loader.fetchFinalScoresByEventsIDs(ctx, eventIDs)
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

	for _, s := range scores {
		if e := lookup[s.EventID]; e != nil && e.Status != models.Scheduled {
			e.FinalScores = append(e.FinalScores, s.FinalScore)
		}
	}

	return events, nil
}

func (s *EventService) insertEventTX(ctx context.Context, req models.CreateEventReq) (*models.DetailedEvent, error) {
	var eventParams repo.InsertEventParams
	if err := copier.Copy(&eventParams, &req); err != nil {
		return nil, err
	}

	startTime, _ := time.Parse(util.TimestampLayout, req.StartTime)
	eventParams.StartTime = startTime
	eventParams.Status = repo.StatusScheduled

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.q.WithTx(tx)

	eventID, err := qtx.InsertEvent(ctx, eventParams)
	if err != nil {
		return nil, err
	}

	participantsParams := repo.InsertParticipantsParams{
		EventID:  eventID,
		TeamsIds: req.TeamIDs,
	}

	if err := qtx.InsertParticipants(ctx, participantsParams); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.GetEvent(ctx, eventID)
}
