package models

type Player struct {
	PlayerID    int32   `json:"player_id,omitempty"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	ShirtNumber int32   `json:"shirt_number"`
	Country     Country `json:"country"`
}

type DetailedPlayer struct {
	Player
	SportName string `json:"sport_name"`
	Height    int32  `json:"height"`
	Age       int32  `json:"age"`
	Photo     string `json:"photo"`
}
