# ==========================================
# Makefile for k8s-devguardian-ai
# AI-powered Kubernetes security auditing tool
# ==========================================
#
# This Makefile provides various commands for building, testing, and running
# the k8s-devguardian-ai tool. It supports multiple platforms and architectures,
# as well as Docker containerization.
#
# For a complete list of available commands, run: make help

# ==========================================
# Go parameters
# ==========================================
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

# ==========================================
# Application parameters
# ==========================================
BINARY_NAME=devguardian
VERSION=0.1.0
BUILD_TIME=$(shell date +%FT%T%z)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
MAIN_PATH=./main.go

# ==========================================
# Build flags
# ==========================================
LDFLAGS=-ldflags "-w -s -X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.GitBranch=$(GIT_BRANCH)'"

# ==========================================
# Output directories
# ==========================================
BIN_DIR=bin
DIST_DIR=dist
TEST_OUTPUT_DIR=test-output

# ==========================================
# Platform targets
# ==========================================
PLATFORMS=linux darwin windows
ARCHS=amd64 arm64

# ==========================================
# Define phony targets (targets that don't create files)
# ==========================================
.PHONY: all build build-all clean test coverage deps tidy fmt lint vet run install help docker docker-run release version check

# ==========================================
# Main targets
# ==========================================

# Default target: show help
.DEFAULT_GOAL := help

# Build and test everything
all: clean deps tidy fmt test build

# Build the application
# Usage examples:
#   make build                      - Build with default settings
#   make build VERSION=1.0.0       - Build with custom version
#   make build BUILDFLAGS="-race"   - Build with race detector
#   make build OUTPUT=./devguardian - Build to custom output path
build: prepare
	@echo "🔨 Building $(BINARY_NAME) v$(VERSION)..."
	if [ -z "$(OUTPUT)" ]; then \
		$(GOBUILD) $(BUILDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) $(LDFLAGS) $(MAIN_PATH); \
		@echo "✅ Build complete: $(BIN_DIR)/$(BINARY_NAME)"; \
	else \
		$(GOBUILD) $(BUILDFLAGS) -o $(OUTPUT) $(LDFLAGS) $(MAIN_PATH); \
		@echo "✅ Build complete: $(OUTPUT)"; \
	fi

# Build for all platforms
build-all: prepare
	@echo "🔨 Building $(BINARY_NAME) for all platforms..."
	$(foreach GOOS, $(PLATFORMS),\
		$(foreach GOARCH, $(ARCHS),\
			$(shell GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH) $(LDFLAGS) $(MAIN_PATH) 2>/dev/null && echo "✅ Built for $(GOOS)/$(GOARCH)" || echo "❌ Failed to build for $(GOOS)/$(GOARCH)")\
		)\
	)

# Create necessary directories
prepare:
	@mkdir -p $(BIN_DIR) $(DIST_DIR) $(TEST_OUTPUT_DIR)

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	$(GOCLEAN)
	rm -rf $(BIN_DIR) $(DIST_DIR) $(TEST_OUTPUT_DIR)
	rm -f coverage.out
	@echo "✅ Clean complete"

# ==========================================
# Development targets
# ==========================================

# Run tests
# Usage examples:
#   make test                                      - Run all tests
#   make test PKG="./internal/ai"                  - Run tests in the ai package
#   make test PKG="./internal/output"              - Run tests in the output package
#   make test PKG="./internal/ai ./internal/output" - Run tests in multiple packages
#   make test TESTFLAGS="-v -race"                 - Run tests with custom flags
test: prepare
	@echo "🧪 Running tests..."
	if [ -z "$(PKG)" ]; then \
		$(GOTEST) $(TESTFLAGS) ./... -coverprofile=$(TEST_OUTPUT_DIR)/coverage.out; \
	else \
		$(GOTEST) $(TESTFLAGS) $(PKG) -coverprofile=$(TEST_OUTPUT_DIR)/coverage.out; \
	fi
	@echo "✅ Tests complete"

# Generate test coverage report
coverage: test
	@echo "📊 Generating coverage report..."
	$(GOCMD) tool cover -html=$(TEST_OUTPUT_DIR)/coverage.out -o $(TEST_OUTPUT_DIR)/coverage.html
	@echo "✅ Coverage report generated: $(TEST_OUTPUT_DIR)/coverage.html"

# Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	$(GOGET) -v ./...
	@echo "✅ Dependencies downloaded"

# Tidy Go modules
tidy:
	@echo "🧹 Tidying Go modules..."
	$(GOMOD) tidy
	@echo "✅ Modules tidied"

# Format code
fmt:
	@echo "✨ Formatting code..."
	$(GOCMD) fmt ./...
	@echo "✅ Code formatted"

# Lint code
lint:
	@echo "🔍 Linting code..."
	$(GOLINT) run
	@echo "✅ Linting complete"

# Vet code
vet:
	@echo "🔍 Vetting code..."
	$(GOVET) ./...
	@echo "✅ Vetting complete"

# Run the application
# Usage examples:
#   make run                                                  - Run basic audit
#   make run ARGS="audit --output json"                       - Run audit with JSON output
#   make run ARGS="audit --output html --file report.html"    - Run audit with HTML output to file
#   make run ARGS="audit --ai-provider openai --api-key KEY"  - Run audit with OpenAI
#   make run ARGS="audit --ai-provider ollama"                - Run audit with Ollama
#   make run ARGS="--help"                                    - Show help
#   make run ARGS="audit --help"                              - Show audit command help
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	if [ -z "$(ARGS)" ]; then \
		./$(BIN_DIR)/$(BINARY_NAME) audit; \
	else \
		./$(BIN_DIR)/$(BINARY_NAME) $(ARGS); \
	fi

# Install the application
install: build
	@echo "📦 Installing $(BINARY_NAME) to /usr/local/bin/..."
	cp $(BIN_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ Installation complete"

# ==========================================
# Docker targets
# ==========================================

# Build Docker image
# Usage examples:
#   make docker                      - Build with default tag (version from Makefile)
#   make docker VERSION=latest       - Build with custom tag
#   make docker DOCKER_ARGS="--no-cache" - Build with additional Docker arguments
docker:
	@echo "🐳 Building Docker image..."
	docker build $(DOCKER_ARGS) -t $(BINARY_NAME):$(VERSION) .
	@echo "✅ Docker image built: $(BINARY_NAME):$(VERSION)"

# Run Docker container
# Usage examples:
#   make docker-run                  - Run container with default settings
#   make docker-run VERSION=latest   - Run container with specific tag
#   make docker-run DOCKER_RUN_ARGS="-e API_KEY=xyz" - Run with environment variables
#   make docker-run CMD="audit --output json"        - Run with custom command
docker-run: docker
	@echo "🐳 Running Docker container..."
	if [ -z "$(CMD)" ]; then \
		docker run --rm -it $(DOCKER_RUN_ARGS) $(BINARY_NAME):$(VERSION); \
	else \
		docker run --rm -it $(DOCKER_RUN_ARGS) $(BINARY_NAME):$(VERSION) $(CMD); \
	fi

# ==========================================
# Help target
# ==========================================

# Show help
help:
	@echo "🛠️  k8s-devguardian-ai Makefile Help 🛠️"
	@echo "=========================================="
	@echo "Main Targets:"
	@echo "  help        - Show this help message"
	@echo "  all         - Clean, download deps, tidy, format, test, and build"
	@echo "  build       - Build the application"
	@echo "  build-all   - Build for all platforms (linux, darwin, windows)"
	@echo "  clean       - Clean build artifacts"
	@echo "  run         - Build and run the application"
	@echo "  install     - Install binary to /usr/local/bin/"
	@echo ""
	@echo "Development Targets:"
	@echo "  test        - Run tests"
	@echo "  coverage    - Run tests with coverage report"
	@echo "  deps        - Download dependencies"
	@echo "  tidy        - Tidy Go modules"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  vet         - Vet code"
	@echo ""
	@echo "Docker Targets:"
	@echo "  docker      - Build Docker image"
	@echo "  docker-run  - Run Docker container"
	@echo ""
	@echo "Additional Targets:"
	@echo "  version     - Show version information"
	@echo "  release     - Create release archives for all platforms"
	@echo "  check       - Run all code quality checks (fmt, lint, vet)"
	@echo ""
	@echo "Command Combinations:"
	@echo "  make                    - Show help (default target)"
	@echo "  make build run          - Build and run the application"
	@echo "  make clean build        - Clean and rebuild the application"
	@echo "  make fmt lint vet       - Format, lint, and vet the code"
	@echo "  make test coverage      - Run tests and generate coverage report"
	@echo "  make clean deps tidy    - Clean, download deps, and tidy modules"
	@echo "  make build install      - Build and install the application"
	@echo "  make docker docker-run  - Build and run Docker container"
	@echo "  make all                - Run clean, deps, tidy, fmt, test, build"
	@echo ""
	@echo "Advanced Combinations:"
	@echo "  make clean deps tidy fmt lint vet test build      - Full development cycle"
	@echo "  make clean build-all                             - Rebuild for all platforms"
	@echo "  make clean deps tidy fmt test docker docker-run  - Test and run in Docker"
	@echo "  make clean build run                            - Clean, build, and run"
	@echo "=========================================="
	@echo "Version: $(VERSION) ($(GIT_BRANCH):$(GIT_COMMIT))"
	@echo "Build time: $(BUILD_TIME)"
	@echo "=========================================="

# ==========================================
# Additional targets
# ==========================================

# Show version information
version:
	@echo "$(BINARY_NAME) version $(VERSION)"
	@echo "Git commit: $(GIT_COMMIT)"
	@echo "Git branch: $(GIT_BRANCH)"
	@echo "Build time: $(BUILD_TIME)"

# Create a release (builds for all platforms and creates archives)
# Usage: make release VERSION=1.0.0
release: clean build-all
	@echo "📦 Creating release v$(VERSION)..."
	mkdir -p $(DIST_DIR)/release
	$(foreach GOOS, $(PLATFORMS),\
		$(foreach GOARCH, $(ARCHS),\
			if [ -f "$(DIST_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)" ]; then \
				tar -czf $(DIST_DIR)/release/$(BINARY_NAME)-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz -C $(DIST_DIR) $(BINARY_NAME)-$(GOOS)-$(GOARCH); \
				echo "✅ Created archive for $(GOOS)/$(GOARCH)"; \
			fi; \
		)\
	)
	@echo "✅ Release v$(VERSION) created in $(DIST_DIR)/release"

# Check code quality and potential issues
check: fmt lint vet
	@echo "🔍 Running static analysis..."
	$(GOVET) -shadow ./...
	@echo "✅ Static analysis complete"

# ==========================================
# Command examples (documented in comments)
# ==========================================
#
# Basic development workflow:
#   make clean deps tidy fmt lint vet test build
#
# Quick build and run:
#   make build run
#
# Build with custom version and run with JSON output:
#   make build VERSION=1.0.0 run ARGS="audit --output json"
#
# Build for all platforms and create a release:
#   make release VERSION=1.0.0
#
# Run tests for specific packages with race detection:
#   make test PKG="./internal/ai ./internal/output" TESTFLAGS="-race"
#
# Build and run in Docker with custom command:
#   make docker-run CMD="audit --output html --file /tmp/report.html"
#
# Install locally:
#   make install
#
# Show version information:
#   make version
