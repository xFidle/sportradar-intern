package models

type FinalScore struct {
	TeamID   int32 `json:"team_id"`
	AggScore int32 `json:"agg_score"`
}

type Score struct {
	TeamID  int32 `json:"team_id"`
	Segment int32 `json:"segment_id"`
	Score   int32 `json:"socre"`
}
