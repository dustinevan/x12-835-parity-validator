IMAGE ?= x12-835-validator:local
PORT ?= 8089
DOCKER ?= docker
COMPOSE_PROJECT_NAME ?= x12-835-validator
APP_SERVER_WEB_API_HOST ?= 127.0.0.1:$(PORT)
# Where make install.<name> puts a program: go env GOBIN, or GOPATH/bin when
# GOBIN is unset (the same place go install uses).
GOBIN := $(or $(shell go env GOBIN),$(shell go env GOPATH)/bin)

export IMAGE PORT COMPOSE_PROJECT_NAME

# Colors for list and help (tput; blank when there is no terminal).
YELLOW := $(shell tput -Txterm setaf 3 2>/dev/null)
CYAN := $(shell tput -Txterm setaf 6 2>/dev/null)
RESET := $(shell tput -Txterm sgr0 2>/dev/null)

.DEFAULT_GOAL := help

# Every target carries a "## description" on its line; help reads them, so
# the list never goes stale. list is the same commands grouped.

help: ## Every target with its description, read from the comment beside it
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_.-]+:.*## / {printf "  ${YELLOW}%-18s${RESET} %s\n", $$1, $$2}' $(MAKEFILE_LIST)

list: ## The commands, grouped
	@echo "${CYAN}X12 835 Parity Validator${RESET}"
	@echo ""
	@echo "${YELLOW}Build:${RESET}"
	@echo "  make vdb          The batch tool -> bin/vdb (migrate, load, compare, summarize, seed)"
	@echo "  make check        The gate: gofmt, go vet, go test, frontend typecheck and build"
	@echo ""
	@echo "${YELLOW}Run locally (one process per terminal; Ctrl+C stops it):${RESET}"
	@echo "  make app          The Go server on http://localhost:$(PORT), forwarding the UI to Vite; run make ui beside it"
	@echo "  make ui           The app's UI with hot reload (Vite on 5179; open the app at $(PORT), not 5179)"
	@echo "  make docs         The docs site on http://localhost:5180 (local only)"
	@echo "  make demo         The app as one binary with the UI bundled in, for demos and the tunnel"
	@echo "  make tunnel       A public Cloudflare link to localhost:$(PORT); the URL prints in the output"
	@echo "  make db           Create, migrate, and seed the local Postgres database"
	@echo "                    (bin/vdb migrate down|up -local -yes and seed -local do the steps one at a time)"
	@echo ""
	@echo "${YELLOW}Docker:${RESET}"
	@echo "  make image        Build the image: /validator and /vdb, no sample data"
	@echo "  make docker-up    Build and start Compose in the background"
	@echo "  make docker-stop  Stop Compose; keep its data"
	@echo "  make docker-logs  Follow the container logs"
	@echo "  make docker-status  Show container status"
	@echo ""
	@echo "${YELLOW}Data and code generation:${RESET}"
	@echo "  make gen          Synthetic MAC and vendor 835s -> data/generated"
	@echo "  make codes        Update the CARC and RARC entries tables from the JSON lists"
	@echo "  make tidy         go mod tidy and refresh vendor/ after a go.mod change"
	@echo ""
	@echo "${YELLOW}Tools (installed into $(GOBIN), so each runs from any directory):${RESET}"
	@echo "  make install.tools      All four below, trying each and reporting the failures at the end"
	@echo "  make install.vdb        vdb: the schema migrations and the batch pipeline"
	@echo "  make install.validator  validator: the service with the UI embedded (builds the UI first)"
	@echo "  make install.gen        gen: synthetic 835s and the truth file"
	@echo "  make install.codegen    codegen: the CARC and RARC entries tables from their JSON lists"
	@echo ""
	@echo "PORT=8090 changes the local, tunnel, and Docker port (default $(PORT))."

vdb: ## The batch tool -> bin/vdb (migrate, load, compare, summarize, seed)
	go build -o bin/vdb ./backend/app/vdb

# Local Postgres for the views (Postgres.app by default): create the validator
# database if missing, apply migrations, and load and compare the generated
# test data. make app needs this database; see README "Run with the database".
PG_BIN ?= /Applications/Postgres.app/Contents/Versions/latest/bin

db: vdb ## Create, migrate, and seed the local Postgres database
	$(PG_BIN)/psql -h localhost -U postgres -Atc "SELECT 1 FROM pg_database WHERE datname = 'validator'" | grep -q 1 || $(PG_BIN)/createdb -h localhost -U postgres validator
	./bin/vdb migrate up -local -yes
	test -d data/generated || $(MAKE) gen
	./bin/vdb seed -local

demo: bundle ## The app on localhost:8089 (or PORT) as one binary with the UI bundled in, for demos and the tunnel
	go build -o bin/validator ./backend/app
	APP_SERVER_WEB_API_HOST="$(APP_SERVER_WEB_API_HOST)" ./bin/validator

frontend/node_modules/.package-lock.json: frontend/package.json frontend/package-lock.json
	npm --prefix frontend ci

bundle: frontend/node_modules/.package-lock.json ## Build the app's UI into the bundle the binary embeds (check, demo, and the image use it)
	cd frontend && npm run build

# The binary embeds backend/app/site, so the directory must exist before Go
# can build. The first make app builds the bundle once; after that the app
# is served from Vite and the embedded copy is only a placeholder.
backend/app/site/dist/index.html: frontend/node_modules/.package-lock.json
	cd frontend && npm run build

# Development is three terminals, one process each, Ctrl+C stops each one:
# make app (the Go server on 8089, forwarding UI requests to Vite on 5179),
# make ui (Vite for the app's UI, hot reload), make docs (the docs site on
# 5180). The local database must be up (make db). Open the app at
# http://localhost:8089, never at 5179.
app: backend/app/site/dist/index.html ## The Go server on localhost:8089 (or PORT), forwarding the UI to Vite on 5179; run make ui beside it
	go build -o bin/validator ./backend/app
	APP_SERVER_WEB_UI_PROXY=http://localhost:5179 APP_SERVER_WEB_API_HOST="$(APP_SERVER_WEB_API_HOST)" ./bin/validator

ui: frontend/node_modules/.package-lock.json ## The app's UI with hot reload, Vite on 5179 (make app serves it on 8089)
	cd frontend && npx vite --port 5179

docs: frontend/node_modules/.package-lock.json ## The docs site, Vite on localhost:5180, local only
	cd frontend && npx vite --config vite.docs.config.ts

# A Cloudflare quick tunnel to the local app: a temporary public URL,
# printed in the output, that dies with this process. For demos of
# synthetic data only.
tunnel: ## A public Cloudflare link to localhost:8089 (or PORT); the URL prints in the output
	cloudflared tunnel --no-autoupdate --url http://localhost:$(PORT)

image: ## Build the Docker image: /validator and /vdb, no sample data
	$(DOCKER) build -t "$(IMAGE)" .

docker-up: ## Build and start Compose in the background
	$(DOCKER) compose up --build -d

docker-stop: ## Stop Compose; keep its data
	$(DOCKER) compose stop

docker-logs: ## Follow the container logs
	$(DOCKER) compose logs --follow

docker-status: ## Show container status
	$(DOCKER) compose ps

# Synthetic MAC and vendor 835s plus truth.jsonl (tools/gen, its own module).
gen: ## Synthetic MAC and vendor 835s -> data/generated
	cd tools/gen && go run . -out ../../data/generated

# Rewrite the rows of the entries tables in backend/pkg/carc/carc.go and
# backend/pkg/rarc/rarc.go from carc.json/rarc.json; nothing else in those
# files changes.
# Run after refreshing either JSON file. codegen runs from its own directory
# and is told where the repository root is.
codes: ## Update the CARC and RARC entries tables from the JSON lists
	cd tools/codegen && go run . -root ../..

# Each program installs into GOBIN under its own name, so it is typed from
# any directory.
install.vdb: ## Install vdb (the schema migrations and the batch pipeline) into GOBIN
	@go build -o $(GOBIN)/vdb ./backend/app/vdb && echo "vdb installed in $(GOBIN)"

install.validator: bundle ## Install validator (the service, UI embedded) into GOBIN
	@go build -o $(GOBIN)/validator ./backend/app && echo "validator installed in $(GOBIN)"

install.gen: ## Install gen (synthetic 835s and the truth file) into GOBIN
	@cd tools/gen && go build -o $(GOBIN)/gen . && echo "gen installed in $(GOBIN)"

install.codegen: ## Install codegen (the CARC and RARC entries tables from their JSON lists) into GOBIN
	@cd tools/codegen && go build -o $(GOBIN)/codegen . && echo "codegen installed in $(GOBIN)"

# Tries every program and reports the failures at the end rather than
# stopping at the first.
install.tools: ## Install vdb, validator, gen, and codegen into GOBIN; failures are reported at the end
	@failed=""; \
	for t in vdb validator gen codegen; do \
		$(MAKE) --no-print-directory install.$$t || failed="$$failed $$t"; \
	done; \
	if [ -n "$$failed" ]; then echo "install.tools: failed:$$failed"; exit 1; fi; \
	echo "install.tools: vdb, validator, gen, and codegen installed in $(GOBIN)"

# Go dependencies are vendored (vendor/), so builds need no module download.
# Run after any change to go.mod. tools/gen uses only the standard library
# and has nothing to vendor.
tidy: ## go mod tidy and refresh vendor/ after a go.mod change
	go mod tidy && go mod vendor

# gofmt reads this module's own packages only, not vendor/.
check: bundle ## The gate: gofmt, go vet, go test, frontend typecheck and build
	gofmt -l $$(go list -f '{{.Dir}}' ./...) && go vet ./... && go test ./...

.PHONY: help list vdb db app ui docs demo bundle tunnel image docker-up docker-stop docker-logs docker-status gen codes install.vdb install.validator install.gen install.codegen install.tools tidy check
