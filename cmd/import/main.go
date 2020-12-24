package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/elastic/go-elasticsearch/v7"

	"github.com/evertras/suumo-scraper/internal/esdata"
	"github.com/evertras/suumo-scraper/internal/suumo"
)

func main() {
	///////////////////////////////////////////////////////////////////////////
	// Prepare ES
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	log.Println(res)
	res.Body.Close()

	///////////////////////////////////////////////////////////////////////////
	// Do the importing
	importer := esdata.NewImporter(es)

	err = importer.DeleteAllDataIndices()

	if err != nil {
		log.Fatal("importer.DeleteAllDataIndices:", err)
	}

	ctx := context.TODO()

	err = filepath.Walk("./data/listings", func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		log.Println("Loading", p)
		f, err := os.Open(p)

		if err != nil {
			return fmt.Errorf("os.Open: %w", err)
		}

		decoder := json.NewDecoder(f)

		var listings []suumo.Listing
		decoder.Decode(&listings)
		err = importer.ImportListings(ctx, listings)

		if err != nil {
			return fmt.Errorf("importer.ImportListings: %w", err)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
