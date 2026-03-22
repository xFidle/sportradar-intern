package models

type Team struct {
	TeamID       int    `json:"team_id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Logo         string `json:"logo"`
}

type DetailedTeam struct {
	Team
	City    City     `json:"city"`
	Players []Player `json:"players"`
}
