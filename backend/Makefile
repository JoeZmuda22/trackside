.PHONY: build run test clean dev seed lint

# Build the server binary
build:
	CGO_ENABLED=1 go build -o bin/server ./cmd/server/

# Run the server
run: build
	./bin/server

# Run the server with seed data
seed: build
	./bin/server --seed

# Run with hot-reload using go run
dev:
	go run ./cmd/server/

# Run with seed data and hot-reload
dev-seed:
	go run ./cmd/server/ --seed

# Run all tests
test:
	go test ./tests/ -v -count=1 -timeout=60s

# Run tests with coverage
test-cover:
	go test ./tests/ -v -count=1 -timeout=60s -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html trackside.db uploads/

# Lint (requires golangci-lint: brew install golangci-lint)
lint:
	golangci-lint run ./...

# Format code
fmt:
	go fmt ./...

# Download dependencies
deps:
	go mod tidy

# Build Docker image
docker-build:
	docker build -t trackside-backend .

# Run Docker container
docker-run:
	docker run -p 8080:8080 \
		-e JWT_SECRET=your-secret-here \
		-v trackside-data:/app/data \
		-v trackside-uploads:/app/uploads \
		trackside-backend
