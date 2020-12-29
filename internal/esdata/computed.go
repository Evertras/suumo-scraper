package esdata

import "github.com/evertras/suumo-scraper/internal/suumo"

type computed struct {
	YenPerSquareMeter float32 `json:"yenPerSquareMeter"`
}

type entry struct {
	suumo.Listing
	computed
}

func toEntry(listing suumo.Listing) entry {
	return entry {
		Listing: listing,
		computed: computed {
			YenPerSquareMeter: float32(listing.PricePerMonthYen) / listing.SquareMeters,
		},
	}
}
