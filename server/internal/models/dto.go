package models

type Filter struct {
	StartAfter    string  `json:"start_after"              validate:"required"`
	EndBefore     string  `json:"end_before"               validate:"required"`
	SportID       *int32  `json:"sport_id,omitempty"`
	CompetitionID *int32  `json:"competition_id,omitempty"`
	TeamIDs       []int32 `json:"team_ids,omitempty"`
}

type CreateEventReq struct {
	CompetitionID int32   `json:"competition_id" validate:"required"`
	VenueID       int32   `json:"venue_id"       validate:"required"`
	StageID       int32   `json:"stage_id"       validate:"required"`
	StartTime     string  `json:"start_time"     validate:"required"`
	TeamIDs       []int32 `json:"team_ids"       validate:"required"`
}
