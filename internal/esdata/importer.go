package esdata

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/evertras/suumo-scraper/internal/suumo"
)

type Importer struct {
	esClient *elasticsearch.Client
}

func NewImporter(esClient *elasticsearch.Client) *Importer {
	return &Importer{
		esClient: esClient,
	}
}

func (i *Importer) DeleteAllDataIndices() error {
	toDelete := []string{
		IndexNameListing,
	}

	for _, index := range toDelete {
		// Do these separately since they may not all exist
		res, err := i.esClient.Indices.Delete([]string{index})

		if err != nil {
			return fmt.Errorf("i.esClient.Indices.Delete %q: %w", index, err)
		}

		if res.StatusCode/200 != 1 && res.StatusCode != 404 {
			return fmt.Errorf("unexpected status code for %q: %d", index, res.StatusCode)
		}
	}

	return nil
}

func (i *Importer) ImportListings(ctx context.Context, listings []suumo.Listing) error {
	err := i.prepLocationMapping(IndexNameListing, []string{"location"})
	if err != nil {
		return fmt.Errorf("i.prepLocationMapping: %w", err)
	}

	bulk, err := startBulkAdder(ctx, i.esClient, IndexNameListing)
	if err != nil {
		return fmt.Errorf("startBulkAdder: %w", err)
	}
	defer bulk.closeWithLoggedError(ctx)

	for i, entry := range listings {
		err = bulk.add(ctx, entry)
		if err != nil {
			return fmt.Errorf("bulk.add #%d: %w", i, err)
		}
	}

	return nil
}

func (i *Importer) prepLocationMapping(index string, locFields []string) error {
	exists, err := i.esClient.Indices.Exists([]string{index})

	if exists.StatusCode == 200 {
		return nil
	}

	_, err = i.esClient.Indices.Create(index)

	if err != nil {
		return fmt.Errorf("esapi.IndicesCreate: %w", err)
	}

	mapping := make(map[string]string)
	for _, field := range locFields {
		mapping[field] = "geo_point"
	}
	indexMappingBody, err := genIndexMappingBody(mapping)

	res, err := i.esClient.Indices.PutMapping(
		indexMappingBody,
		i.esClient.Indices.PutMapping.WithIndex(index),
	)

	if err != nil {
		return fmt.Errorf("creating index with mapping: %w", err)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected error code %d", res.StatusCode)
	}

	return nil
}
