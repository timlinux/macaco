.PHONY: build clean test fmt lint run server client nix-build vendor deps install help start-backend stop-backend status

# Default target
all: build

# Build the application
build:
	go build -o bin/macaco ./cmd/macaco

# Build with version info
build-release:
	go build -ldflags="-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/macaco ./cmd/macaco

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f result result-*

# Run tests
test:
	go test ./... -v

# Run tests with coverage
test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Vet code
vet:
	go vet ./...

# Run in combined mode (default)
run: build
	./bin/macaco

# Run with specific round type
run-beginner: build
	./bin/macaco --round beginner

run-intermediate: build
	./bin/macaco --round intermediate

run-advanced: build
	./bin/macaco --round advanced

run-expert: build
	./bin/macaco --round expert

# Run server only
server: build
	./bin/macaco --server

# Run client only
client: build
	./bin/macaco --client

# Build with Nix
nix-build:
	nix build

# Vendor dependencies
vendor:
	go mod vendor

# Tidy dependencies
deps:
	go mod tidy

# Install to GOPATH
install:
	go install ./cmd/macaco

# Backend management
start-backend: build
	./scripts/start-backend.sh

stop-backend:
	./scripts/stop-backend.sh

status:
	./scripts/status-backend.sh

# Documentation
docs:
	mkdocs serve

docs-build:
	mkdocs build

# Development
dev: build
	./bin/macaco

# Help
help:
	@echo "MoCaCo - Motion Capture Combatant"
	@echo ""
	@echo "Build targets:"
	@echo "  build          Build the application"
	@echo "  build-release  Build with version info"
	@echo "  clean          Remove build artifacts"
	@echo "  nix-build      Build using Nix"
	@echo ""
	@echo "Development targets:"
	@echo "  run            Run in combined mode"
	@echo "  server         Run server only"
	@echo "  client         Run client only"
	@echo "  test           Run tests"
	@echo "  test-cover     Run tests with coverage"
	@echo "  fmt            Format code"
	@echo "  lint           Run linter"
	@echo "  vet            Run go vet"
	@echo ""
	@echo "Backend management:"
	@echo "  start-backend  Start backend server"
	@echo "  stop-backend   Stop backend server"
	@echo "  status         Check backend status"
	@echo ""
	@echo "Documentation:"
	@echo "  docs           Serve documentation locally"
	@echo "  docs-build     Build documentation"
	@echo ""
	@echo "Dependencies:"
	@echo "  deps           Tidy go.mod"
	@echo "  vendor         Vendor dependencies"
	@echo "  install        Install to GOPATH"
