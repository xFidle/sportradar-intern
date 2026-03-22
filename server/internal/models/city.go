package models

type City struct {
	Name    string  `json:"name"`
	Country Country `json:"country"`
}
