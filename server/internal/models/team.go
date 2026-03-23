package models

type Team struct {
	TeamID       int32  `json:"team_id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	LogoPath     string `json:"logo"         copier:"-"`
}

type DetailedTeam struct {
	Team
	City    City     `json:"city"`
	Players []Player `json:"players"`
}
