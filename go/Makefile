SHELL := /bin/bash

.PHONY: test vet fmt fmt-check ci

test:
	go test ./...

vet:
	go vet ./...

# Format all non-empty, non-underscored Go files
fmt:
	find . -type f -name '*.go' ! -name '_*.go' ! -size 0 -print0 | xargs -0 gofmt -s -w

# Check formatting; fails if any files need changes
fmt-check:
	@out=$$(find . -type f -name '*.go' ! -name '_*.go' ! -size 0 -print0 | xargs -0 gofmt -s -l); \
	if [[ -n "$$out" ]]; then echo "The following files need gofmt:"; echo "$$out"; exit 1; fi

# Run basic CI checks locally
ci: vet fmt-check test

.PHONY: db-migrate
db-migrate:
	@[ -n "$$DATABASE_URL" ] || (echo "DATABASE_URL not set" && exit 1)
	go run ./cmd/migrate --dir _data/_db --database "$$DATABASE_URL"

.PHONY: build-default build-otel build-pulsar build-all
build-default:
	go build ./...

build-otel:
	go build -tags=otel ./...

build-pulsar:
	go build -tags=pulsar ./...

build-all:
	go build -tags="otel pulsar" ./...

.PHONY: run-api
run-api:
	API_PORT=$${API_PORT:-8080} go run ./cmd/api
