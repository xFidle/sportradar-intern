package service

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

	event.Competition = models.Competition{
		Name: row.CompetitionName,
		Type: models.CompetitionType(row.CompetitionType),
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
	rows, err := l.q.ListPlayersByTeamIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	players := make([]teamPlayer, 0, 11)
	for _, row := range rows {
		var player models.Player
		if err := copier.Copy(&player, &row); err != nil {
			return nil, err
		}
		player.Country = models.Country{Name: row.CountryName, Code: row.CountryCode}
		players = append(players, teamPlayer{teamID: row.TeamID, player: player})
	}

	return players, nil
}

func (l *loader) fetchEventsByFilter(ctx context.Context, filter models.Filter) ([]models.Event, error) {
	params := repo.ListEventsByFilterParams{}
	if err := copier.Copy(&params, &filter); err != nil {
		return nil, err
	}

	if filter.Status != nil {
		params.Status = repo.NullStatus{
			Status: repo.Status(string(*filter.Status)),
			Valid:  true,
		}
	}

	after, _ := time.Parse(util.DateLayout, filter.StartAfter)
	before, _ := time.Parse(util.DateLayout, filter.EndBefore)

	params.StartAfter = after
	params.EndBefore = before

	rows, err := l.q.ListEventsByFilter(ctx, params)
	if err != nil {
		return nil, err
	}

	events := make([]models.Event, len(rows))
	for i := range events {
		if err := copier.Copy(&events[i], &rows[i]); err != nil {
			return nil, err
		}
		events[i].Competition = models.Competition{
			Name: rows[i].CompetitionName,
			Type: models.CompetitionType(rows[i].CompetitionType),
		}

		if rows[i].CompetitionLogo != nil {
			events[i].Competition.LogoPath = fmt.Sprintf("%s/%s", l.fAddr, *rows[i].CompetitionLogo)
		}
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

func (l *loader) fetchScoresByEventID(ctx context.Context, id int32) ([]models.Score, error) {
	rows, err := l.q.ListScoresByEventID(ctx, id)
	if err != nil {
		return nil, err
	}

	var scores []models.Score
	if err := copier.Copy(&scores, &rows); err != nil {
		return nil, err
	}

	return scores, nil
}

type eventFinalScore struct {
	EventID    int32
	FinalScore models.FinalScore
}

func (l *loader) fetchFinalScoresByEventsIDs(ctx context.Context, eventIDs []int32) ([]eventFinalScore, error) {
	rows, err := l.q.ListFinalScoresByEventsIDs(ctx, eventIDs)
	if err != nil {
		return nil, err
	}

	finalScores := make([]eventFinalScore, 0, len(rows))
	for _, row := range rows {
		var finalScore models.FinalScore
		if err := copier.Copy(&finalScore, &row); err != nil {
			return nil, err
		}
		finalScores = append(finalScores, eventFinalScore{EventID: row.EventID, FinalScore: finalScore})
	}

	return finalScores, nil
}

func (l *loader) checkVenueCorrectnes(ctx context.Context, competitionID, venueID int32) (bool, error) {
	params := repo.IsVenueValidForCompetitionParams{CompetitionID: competitionID, VenueID: venueID}
	return l.q.IsVenueValidForCompetition(ctx, params)
}

func (l *loader) checkTeamsCorrectness(ctx context.Context, competitionID int32, teamIDs []int32) (bool, error) {
	params := repo.CountValidTeamsForCompetitionParams{CompetitionID: competitionID, TeamIds: teamIDs}
	count, err := l.q.CountValidTeamsForCompetition(ctx, params)
	if err != nil {
		return false, err
	}
	return count == int64(len(teamIDs)), err
}
