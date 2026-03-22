package models

type Player struct {
	PlayerID   int32   `json:"player_id"`
	FirstName  string  `json:"first_name"`
	SecondName string  `json:"second_name"`
	Country    Country `json:"country"`
}

type DetailedPlayer struct {
	Player
	SportName string `json:"sport_name"`
	Height    int32  `json:"height"`
	Age       int32  `json:"age"`
	Photo     string `json:"photo"`
}
