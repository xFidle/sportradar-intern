package models

import "time"

type Status string

const (
	Finished   Status = "finished"
	InProgress Status = "in progress"
	Scheduled  Status = "scheduled"
)

type Event struct {
	EventID         int       `json:"event_id"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	Status          Status    `json:"status"`
	SportName       string    `json:"sport_name"`
	CompetitionName string    `json:"competition_name"`
	FinalScore      string    `json:"final_score"`
	Participants    []Team    `json:"participants"`
}

type DetailedEvent struct {
}

type Filter struct {
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	SportID       *int   `json:"sport_id,omitempty"`
	CompetitionID *int   `json:"competition_id,omitempty"`
	TeamID        *int   `json:"team_id,omitempty"`
}
