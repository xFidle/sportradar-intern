package models

import "time"

type Status string

const (
	Finished   Status = "finished"
	InProgress Status = "in progress"
	Scheduled  Status = "scheduled"
)

type baseEvent struct {
	EventID         int32     `json:"event_id"`
	StartTime       time.Time `json:"start_time"`
	Status          Status    `json:"status"`
	SportName       string    `json:"sport_name"`
	CompetitionName string    `json:"competition_name"`
}

type Event struct {
	baseEvent
	Participants []Team `json:"participants"`
}

type DetailedEvent struct {
	baseEvent
	VenueName    string         `json:"venue_name"`
	Scores       []Score        `json:"scores"`
	Participants []DetailedTeam `json:"participants"`
}
