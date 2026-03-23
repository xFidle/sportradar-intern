package event

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/xFidle/sportradar-intern/server/internal/models"
	"github.com/xFidle/sportradar-intern/server/internal/repo"
	"github.com/xFidle/sportradar-intern/server/internal/util"
)

// loader aggregates queries from SQLC and maps models to domain ones
type loader struct {
	fAddr string
	db    *pgxpool.Pool
	q     *repo.Queries
}

func newLoader(db *pgxpool.Pool, fileserverAddr string) *loader {
	return &loader{
		fAddr: fileserverAddr,
		db:    db,
		q:     repo.New(db),
	}
}

func (l *loader) fetchEventByID(ctx context.Context, id int32) (*models.DetailedEvent, error) {
	row, err := l.q.GetDetailedEventByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var event models.DetailedEvent
	if err := copier.Copy(&event, row); err != nil {
		return nil, err
	}

	return &event, nil
}

func (l *loader) fetchTeamsByEventID(ctx context.Context, id int32) ([]models.DetailedTeam, error) {
	result := make([]models.DetailedTeam, 0, 2)
	rows, err := l.q.ListDetailedTeamsByEventID(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		var team models.DetailedTeam
		if err := copier.Copy(&team, &row); err != nil {
			return nil, err
		}

		team.LogoPath = fmt.Sprintf("%s/%s", l.fAddr, team.LogoPath)
		team.City = models.City{
			Name:    row.CityName,
			Country: models.Country{Name: row.CountryName, Code: row.CountryCode},
		}
		result = append(result, team)
	}

	return result, nil
}

type teamPlayer struct {
	teamID int32
	player models.Player
}

func (l *loader) fetchPlayersByTeamID(ctx context.Context, ids []int32) ([]teamPlayer, error) {
	result := make([]teamPlayer, 0, 11)
	rows, err := l.q.ListPlayersByTeamIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		var player models.Player
		if err := copier.Copy(&player, &row); err != nil {
			return nil, err
		}
		player.Country = models.Country{Name: row.CountryName, Code: row.CountryCode}
		result = append(result, teamPlayer{teamID: row.TeamID, player: player})
	}

	return result, nil
}

func (l *loader) fetchEventsByFilter(ctx context.Context, filter models.Filter) ([]models.Event, error) {
	params := repo.ListEventsByFilterParams{}
	if err := copier.Copy(&params, &filter); err != nil {
		return nil, err
	}

	after, _ := time.Parse(util.DateLayout, filter.StartAfter)
	before, _ := time.Parse(util.DateLayout, filter.EndBefore)

	params.StartAfter = after
	params.EndBefore = before

	rows, err := l.q.ListEventsByFilter(ctx, params)
	if err != nil {
		return nil, err
	}

	var events []models.Event
	if err := copier.Copy(&events, &rows); err != nil {
		return nil, err
	}

	return events, nil
}

type eventTeam struct {
	eventID int32
	team    models.Team
}

func (l *loader) fetchTeamsByEventIDs(ctx context.Context, eventIDs []int32) ([]eventTeam, error) {
	result := make([]eventTeam, 0, 2*len(eventIDs))
	rows, err := l.q.ListTeamsByEventsIDs(ctx, eventIDs)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		var team models.Team
		if err := copier.Copy(&team, &row); err != nil {
			return nil, err
		}

		team.LogoPath = fmt.Sprintf("%s/%s", l.fAddr, team.LogoPath)
		result = append(result, eventTeam{eventID: row.EventID, team: team})
	}

	return result, nil
}
