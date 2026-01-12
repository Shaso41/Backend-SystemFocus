.PHONY: all build test bench run clean docker docker-run fmt vet

# Binary name
BINARY_NAME=redis-clone
DOCKER_IMAGE=redis-clone:latest

# Build the application
build:
	@echo "🔨 Building..."
	go build -o $(BINARY_NAME) ./cmd/server

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
coverage: test
	@echo "📊 Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@echo "⚡ Running benchmarks..."
	go test -bench=. -benchmem ./internal/store

# Run the server
run: build
	@echo "🚀 Starting server..."
	./$(BINARY_NAME)

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "🔍 Running go vet..."
	go vet ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	go clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Build Docker image
docker:
	@echo "🐳 Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Run Docker container
docker-run: docker
	@echo "🚀 Running Docker container..."
	docker run -p 6379:6379 --name redis-clone-container $(DOCKER_IMAGE)

# Run with Docker Compose
compose-up:
	@echo "🚀 Starting with Docker Compose..."
	docker-compose up -d

# Stop Docker Compose
compose-down:
	@echo "🛑 Stopping Docker Compose..."
	docker-compose down

# Install dependencies
deps:
	@echo "📦 Downloading dependencies..."
	go mod download

# Run all checks (fmt, vet, test)
check: fmt vet test
	@echo "✅ All checks passed!"

# Default target
all: check build
