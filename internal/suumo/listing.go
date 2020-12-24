package suumo

import "github.com/evertras/suumo-scraper/internal/geocode"

type StationWalk struct {
	Name    string `json:"name"`
	Minutes int    `json:"minutes"`
}

type Listing struct {
	Title            string            `json:"title"`
	Neighborhood     string            `json:"neighborhood"`
	Location         *geocode.Location `json:"location,omitempty"`
	AgeYears         int               `json:"ageYears"`
	Floor            int               `json:"floor"`
	PricePerMonthYen int               `json:"pricePerMonthYen"`
	Layout           string            `json:"layout"`
	SquareMeters     float32           `json:"squareMeters"`
	Ward             Ward              `json:"ward"`
	Direction        string            `json:"direction,omitempty"`
	Stations         []StationWalk     `json:"stations,omitempty"`
}
