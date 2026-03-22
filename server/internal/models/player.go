package models

type Player struct {
	PlayerID   int     `json:"player_id"`
	FirstName  string  `json:"first_name"`
	SecondName string  `json:"second_name"`
	Country    Country `json:"country"`
}

type DetailedPlayer struct {
	Player
	SportName string `json:"sport_name"`
	Height    int    `json:"height"`
	Age       int    `json:"age"`
	Photo     string `json:"photo"`
}
