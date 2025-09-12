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

.PHONY: run tidy air
run:
	go run ./cmd/server

tidy:
	go mod tidy

air:
	air -c air.toml