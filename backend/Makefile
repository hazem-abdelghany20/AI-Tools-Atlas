.PHONY: help run test build clean migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make run          - Run the development server"
	@echo "  make test         - Run all tests"
	@echo "  make build        - Build the application"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make migrate-up   - Run database migrations up"
	@echo "  make migrate-down - Run database migrations down"

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
