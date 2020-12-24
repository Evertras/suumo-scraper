package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/evertras/suumo-scraper/internal/geocode"
	"github.com/evertras/suumo-scraper/internal/suumo"
)

func main() {
	log.Println("Getting wards...")

	wards, err := suumo.ListWards(suumo.PrefectureFromCode(suumo.PrefectureCodeSaitama))

	if err != nil {
		log.Fatal("suumo.ListWards:", err)
	}

	log.Printf("Found %d wards", len(wards))

	wait := sync.WaitGroup{}

	os.MkdirAll("data/listings", os.ModeDir | 0777)

	ctx := context.TODO()

	geocodeClient, err := geocode.NewClient(os.Getenv("GOOGLE_MAPS_API_KEY"))

	for _, ward := range wards {
		wait.Add(1)
		go func(ward suumo.Ward) {
			defer wait.Done()
			path := fmt.Sprintf("data/listings/%s.json", ward.Code)

			if _, err := os.Stat(path); !os.IsNotExist(err) {
				log.Printf("Data already exists for %s (%s), skipping...", ward.Name, ward.Code)
				return
			}

			log.Printf("Fetching %s (%s)", ward.Name, ward.Code)

			listings, err := suumo.WardListings(ctx, geocodeClient, ward)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Total hits:", len(listings))

			file, err := os.Create(path)

			if err != nil {
				log.Fatal("Failed to open file to write:", err)
			}

			encoder := json.NewEncoder(file)

			err = encoder.Encode(listings)

			if err != nil {
				log.Fatal("Failed to encode JSON to file:", err)
			}
		}(ward)
	}

	wait.Wait()

	log.Println("Done!")
}
