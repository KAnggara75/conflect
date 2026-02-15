.PHONY: test coverage coverage-html clean build run help

# Default target
.DEFAULT_GOAL := help

## test: Run all tests
test:
	@echo "Running tests..."
	go test ./... -v

## coverage: Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out -covermode=atomic
	@echo "\n=== Coverage Summary ==="
	go tool cover -func=coverage.out | tail -1

## coverage-html: Generate HTML coverage report
coverage-html: coverage
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## coverage-detail: Show detailed coverage by package
coverage-detail: coverage
	@echo "\n=== Detailed Coverage Report ==="
	go tool cover -func=coverage.out

## test-race: Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	go test ./... -race -v

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test ./... -v -coverprofile=coverage.out
	go tool cover -func=coverage.out

## build: Build the application
build:
	@echo "Building application..."
	go build -o bin/conflect cmd/conflect/conflect.go
	@echo "Build complete: bin/conflect"

## build-prod: Build for production with optimizations
build-prod:
	@echo "Building for production..."
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/conflect cmd/conflect/conflect.go
	@echo "Production build complete: bin/conflect"

## run: Run the application
run:
	@echo "Running application..."
	go run cmd/conflect/conflect.go

## clean: Clean build artifacts and coverage files
clean:
	@echo "Cleaning..."
	rm -f coverage.out coverage.html
	rm -rf bin/
	go clean
	@echo "Clean complete"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

## tidy: Tidy go.mod
tidy:
	@echo "Tidying go.mod..."
	go mod tidy

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## lint: Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: brew install golangci-lint"; \
		exit 1; \
	fi

## check: Run all checks (fmt, vet, test)
check: fmt vet test
	@echo "All checks passed!"

## ci: Run CI pipeline locally
ci: clean deps check coverage
	@echo "CI pipeline complete!"

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
