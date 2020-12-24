default: ensure-environment
	@echo See Makefile for commands

fmt:
	@go fmt ./pkg/...
	@go fmt ./cmd/...

scrape: ensure-environment
	@mkdir -p data
	@go run -race ./cmd/scrape/main.go

import: ensure-environment
	@go run -race ./cmd/import/main.go

geocode: ensure-environment
	@go run -race ./cmd/geocode/main.go

test:
	go test ./internal/...

################################################################################
# Elasticsearch helpers
es:
	@echo "Starting Elasticsearch stack"
	cd ./elastic && docker-compose up

clear-es:
	@echo "Deleting ES indices, some may not be found and that's fine"
	curl -XDELETE 'http://localhost:9200/listing'
	@echo

################################################################################
# Other helpers
ensure-environment:
	@if test -z ${GOOGLE_MAPS_API_KEY}; then echo "Set GOOGLE_MAPS_API_KEY"; exit 1; fi

