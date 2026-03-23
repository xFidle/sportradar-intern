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
	EndTime         time.Time `json:"end_time"`
	Status          Status    `json:"status"`
	SportName       string    `json:"sport_name"`
	CompetitionName string    `json:"competition_name"`
}

type Event struct {
	baseEvent
	FinalScore   string `json:"final_score"`
	Participants []Team `json:"participants"`
}

type DetailedEvent struct {
	baseEvent
	VenueName    string         `json:"venue_name"`
	Scores       []Score        `json:"scores"`
	Participants []DetailedTeam `json:"participants"`
}

type Filter struct {
	StartAfter    string  `json:"start_after"              validate:"required"`
	EndBefore     string  `json:"end_before"               validate:"required"`
	SportID       *int32  `json:"sport_id,omitempty"`
	CompetitionID *int32  `json:"competition_id,omitempty"`
	TeamIDs       []int32 `json:"team_ids,omitempty"`
}
