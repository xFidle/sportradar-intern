package models

type CompetitionType string

const (
	Tournament CompetitionType = "tournament"
	League     CompetitionType = "league"
	Friendly   CompetitionType = "friendly"
	Other      CompetitionType = "other"
)

type Competition struct {
	CompetitionID int32           `json:"competition_id"`
	Name          string          `json:"name"`
	Type          CompetitionType `json:"type"`
	LogoPath      string          `json:"logo_path"`
}
