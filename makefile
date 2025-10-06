.PHONY: localdev_bootstrap
localdev_bootstrap:
	@echo "Starting local development environment"
	docker compose -f docker-compose.yml up -d

# DB_COMMANDS

.PHONY: migration
migration:
	echo "---> Creating a new migration"
	echo "---> Input args... $(filename)"
	@migrate create -ext sql -dir ./db/migrations/ $(filename)

.PHONY: db-migrate
db-migrate:
	@echo "Running db migrations"
	go mod tidy
	@go run db/migrations/main.go up

.PHONY: db-migrate-down
db-migrate-down:
	@echo "Running db migrations down"
	go mod tidy
	@go run db/migrations/main.go down

.PHONY: db-seed
db-seed:
	@echo "Running all db seeds"
	go mod tidy
	@go run db/seeds/main.go

.PHONY: run tidy air
run:
	go run ./cmd/server

tidy:
	go mod tidy

air:
	air -c air.toml

# DEVELOPMENT COMMANDS

.PHONY: deps
deps:
	@echo "Installing dependencies"
	go mod download
	go mod tidy

.PHONY: lint
lint:
	@echo "Running linter"
	golangci-lint run ./cmd/... ./config/... ./db/... ./pkg/... ./utils/...

.PHONY: test
test:
	@echo "Running tests"
	GOFLAGS=-buildvcs=false go test -v -p 1 ./pkg/... ./cmd/... ./config/... ./db/... ./utils/...

.PHONY: test_ci
test_ci:
	@echo "Running tests with coverage"
	GOFLAGS=-buildvcs=false go test -v -p 1 -coverprofile=coverage.out ./pkg/... ./cmd/... ./config/... ./db/... ./utils/...

.PHONY: build
build:
	@echo "Building application"
	mkdir -p build
	GOFLAGS=-buildvcs=false go build -o build/service ./cmd/server