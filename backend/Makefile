.PHONY: help run test build clean migrate-up migrate-down migrate-create

help:
	@echo "Available commands:"
	@echo "  make run                    - Run the development server"
	@echo "  make test                   - Run all tests"
	@echo "  make build                  - Build the application"
	@echo "  make clean                  - Clean build artifacts"
	@echo "  make migrate-up             - Run database migrations up"
	@echo "  make migrate-down           - Run database migrations down"
	@echo "  make migrate-create name=X  - Create new migration file"

run:
	go run cmd/api/main.go

test:
	go test -v ./...

build:
	go build -o bin/api cmd/api/main.go

clean:
	rm -rf bin/

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required. Usage: make migrate-create name=<name>"; \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations -seq $(name)
