package geocode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"googlemaps.github.io/maps"
)

type Client struct {
	googleClient *maps.Client
	cache        map[string]*Location
	mu           sync.Mutex
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func NewClient(googleAPIKey string) (*Client, error) {
	c, err := maps.NewClient(maps.WithAPIKey(googleAPIKey))

	if err != nil {
		return nil, fmt.Errorf("maps.NewClient: %w", err)
	}

	return &Client{
		googleClient: c,
		cache:        make(map[string]*Location),
	}, nil
}

func filenameForNeighborhood(neighborhood string) string {
	return fmt.Sprintf("./data/geocode/%s.json", neighborhood)
}

func (c *Client) tryGetFromFile(neighborhood string) (*Location, error) {
	filename := filenameForNeighborhood(neighborhood)

	file, err := os.Open(filename)

	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	loc := &Location{}

	err = decoder.Decode(loc)

	if err != nil {
		return nil, fmt.Errorf("decoder.Decode: %w", err)
	}

	return loc, nil
}

func (c *Client) GetNeighborhoodLocation(ctx context.Context, neighborhood string) (*Location, error) {
	var loc *Location
	var err error

	// Note that we could probably speed this up on the first run by having mutexes
	// per neighborhood, but let's err on the side of simple/safe for now
	c.mu.Lock()
	defer c.mu.Unlock()

	if loc, ok := c.cache[neighborhood]; ok {
		return loc, nil
	}

	loc, err = c.tryGetFromFile(neighborhood)

	if err != nil {
		return nil, fmt.Errorf("c.tryGetFromFile: %w", err)
	}

	if loc != nil {
		c.cache[neighborhood] = loc
		return loc, nil
	}

	req := &maps.GeocodingRequest{
		Address: neighborhood,
		Region:  "ja",
	}

	res, err := c.googleClient.Geocode(ctx, req)

	if err != nil {
		return nil, fmt.Errorf("googleClient.Geocode: %w", err)
	}

	if len(res) != 1 {
		log.Printf("Unexpected number of responses for neighborhood %q: %d", neighborhood, len(res))
	}

	err = os.MkdirAll("data/geocode", os.ModeDir|0777)

	if err != nil {
		return nil, fmt.Errorf("os.MkdirAll: %w", err)
	}

	outFile, err := os.Create(filenameForNeighborhood(neighborhood))

	if err != nil {
		return nil, fmt.Errorf("os.Create: %w", err)
	}

	defer outFile.Close()

	encoder := json.NewEncoder(outFile)

	// Use the first entry if there are multiple while hoping it's accurate enough...
	loc = &Location{
		Lat: res[0].Geometry.Location.Lat,
		Lon: res[0].Geometry.Location.Lng,
	}

	err = encoder.Encode(loc)

	if err != nil {
		return nil, fmt.Errorf("encoder.Encode: %w", err)
	}

	c.cache[neighborhood] = loc

	return loc, nil
}
