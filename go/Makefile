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

.PHONY: codegen-recurpost codegen-postiz openapi-validate
codegen-recurpost:
	@command -v oapi-codegen >/dev/null 2>&1 || { echo "oapi-codegen not installed. See docs/recurpost/CODEGEN.md"; exit 1; }
	oapi-codegen -generate types,client -package recurpost -o pkg/api/recurpost/generated_oapi.go pkg/api/recurpost/openapi_new.yaml

codegen-postiz:
	@command -v oapi-codegen >/dev/null 2>&1 || { echo "oapi-codegen not installed. See docs/postiz/CODEGEN.md"; exit 1; }
	oapi-codegen -generate types,client -package postiz -o pkg/api/postiz/generated_oapi.go pkg/api/postiz/openapi.yaml

# Validate OpenAPI files (requires Redocly CLI: npm i -g @redocly/cli)
openapi-validate:
	@command -v redocly >/dev/null 2>&1 || { echo "redocly not installed. npm i -g @redocly/cli"; exit 1; }
	redocly lint pkg/api/recurpost/openapi_new.yaml
	redocly lint pkg/api/postiz/openapi.yaml

.PHONY: docs-copy-openapi
docs-copy-openapi:
	@mkdir -p docs/recurpost docs/postiz
	cp pkg/api/recurpost/openapi_new.yaml docs/recurpost/openapi.yaml
	cp pkg/api/postiz/openapi.yaml docs/postiz/openapi.yaml

# --- EOWR (Engine-Oriented Workflow Runtime) ---
.PHONY: eowr-compile eowr-export eowr-compile-all eowr-export-all eowr-bootstrap

eowr-compile:
	@echo "Usage: make eowr-compile IN=pkg/engine/workflows/lead_import.pseudo OUT=pkg/engine/workflows/compiled/lead_import.go" && \
	[ -n "$$IN" ] && [ -n "$$OUT" ] && go run ./cmd/eowr compile "$$IN" "$$OUT"

eowr-export:
	@echo "Usage: make eowr-export IN=pkg/engine/workflows/lead_import.pseudo OUT=pkg/engine/workflows/exports/lead_import.json" && \
	[ -n "$$IN" ] && [ -n "$$OUT" ] && go run ./cmd/eowr export "$$IN" "$$OUT"

eowr-compile-all:
	go run ./cmd/eowr compile-all

eowr-export-all:
	go run ./cmd/eowr export-all

eowr-bootstrap: eowr-compile-all eowr-export-all
