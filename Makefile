.PHONY: help run build install clean test db-create db-migrate

help:
	@echo "Available commands:"
	@echo "  make install     - Install Go dependencies"
	@echo "  make run         - Run the application"
	@echo "  make build       - Build the application binary"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make db-create   - Create database"
	@echo "  make test        - Run tests"

install:
	go mod download

run:
	go run cmd/server/main.go

build:
	go build -o interview-chatbot cmd/server/main.go

clean:
	rm -f interview-chatbot
	rm -rf voice_answers/
	rm -rf tts_cache/

db-create:
	mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS interview_chatbot;"

test:
	go test ./...
