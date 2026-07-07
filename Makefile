BINARY_NAME = git.shi.foo
BUILD_PATH = bin/$(BINARY_NAME)
MAIN_PATH = ./$(BINARY_NAME)
DB_NAME = $(shell sed -n 's/^DSN=.*dbname=\([^ "]*\).*/\1/p' .env 2>/dev/null)
DB_USER = $(shell sed -n 's/^DSN=.*[[:space:]]user=\([^ "]*\).*/\1/p' .env 2>/dev/null)

.PHONY: setup clean tidy build run dev db-reset all

setup:
	@go mod download
	@go mod tidy
	@which air > /dev/null 2>&1 || go install github.com/air-verse/air@latest

clean:
	@rm -rf bin

tidy:
	@go mod tidy

build:
	@go build -o $(BUILD_PATH) $(MAIN_PATH)

run:
	@if [ ! -f $(BUILD_PATH) ]; then $(MAKE) -s build; fi
	@$(BUILD_PATH)

dev:
	@air

db-reset:
	@docker compose exec -T postgres psql -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME) WITH (FORCE);"
	@docker compose exec -T postgres psql -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME);"
	@echo "DB reset: '$(DB_NAME)' dropped and recreated. Restart the app (stop air, then 'make dev') to re-run migrations."

all: setup clean build run

.SILENT:
