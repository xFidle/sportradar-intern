package models

type Score struct {
	TeamID  int32 `json:"team_id"`
	Segment int32 `json:"segment_id"`
	Score   int32 `json:"socre"`
}
