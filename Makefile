SRC_REST := ./cmd/rest/rest.go
BIN_REST := ./bin/rest
BUILD_REST_CMD := go build -o $(BIN_REST) $(SRC_REST)

reload: air
	@GIN_MODE=release air --log.main_only=true --build.cmd "$(BUILD_REST_CMD)" --build.bin "$(BIN_REST)"

run: build
	@GIN_MODE=release $(BIN_REST)

build:
	@$(BUILD_REST_CMD)

migrate:
	@go run ./cmd/migrate/migrate.go

seed:
	@go run ./cmd/seed/seed.go

reset-db: migrate seed

cron:
	@go run ./cmd/schedule/cron.go

air:
	@command -v air > /dev/null || go install github.com/cosmtrek/air@latest

test:
	@go test ./... -v

test-fail:
	@go test ./... -v | fgrep FAIL || echo "No test failed"

cover:
	@go test ./... --cover --coverprofile=cover.out
	@go tool cover -html=cover.out
	@rm cover.out

cover-all:
	@go test ./... --cover --coverprofile=cover.out >> /dev/null
	@go tool cover --func cover.out | grep total
	@rm cover.out

prod:
	@docker compose up rest

dev:
	@docker compose up rest_dev

.PHONY:migration
migration:
	@docker compose up postgis-migration migration

down:
	@docker compose down