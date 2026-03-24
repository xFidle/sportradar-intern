package models

type Venue struct {
	VenueID  int32  `json:"venue_id"`
	City     City   `json:"city"`
	Name     string `json:"name"`
	Capacity int32  `json:"capacity"`
}
