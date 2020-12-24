package main

import (
	"context"
	"log"
	"os"

	"github.com/evertras/suumo-scraper/internal/geocode"
)

func main() {
	ctx := context.TODO()
	c, err := geocode.NewClient(os.Getenv("GOOGLE_MAPS_API_KEY"))

	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	loc, err := c.GetNeighborhoodLocation(ctx, "東京都中央区新富２")

	if err != nil {
		log.Fatal("Failed to get neighborhood:", err)
	}

	log.Println(loc)
}
