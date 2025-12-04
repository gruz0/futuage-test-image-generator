# FutuAge Test Image Generator Makefile
# ======================================

# Variables
BINARY_NAME := futuage-test-image-gen
MAIN_PACKAGE := .
BUILD_DIR := dist
OUTPUT_DIR := test-images
VERSION := $(shell grep -m1 'version = ' cmd/root.go | cut -d'"' -f2)
GO_VERSION := $(shell go version | cut -d' ' -f3)
LDFLAGS := -ldflags="-s -w"
INSTALL_PATH := /usr/local/bin

# Colors for help output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RESET := \033[0m

# Default goal
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo ""
	@echo "$(CYAN)FutuAge Test Image Generator$(RESET) v$(VERSION)"
	@echo "$(YELLOW)Usage:$(RESET) make [target]"
	@echo ""
	@echo "$(YELLOW)Available targets:$(RESET)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

# =============================================================================
# Build Targets
# =============================================================================

.PHONY: build
build: ## Build the binary for current platform
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "✓ Built: ./$(BINARY_NAME)"

.PHONY: build-all
build-all: clean-dist ## Build binaries for all platforms (darwin, linux, windows)
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@echo "  → darwin/arm64..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	@echo "  → darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	@echo "  → linux/amd64..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	@echo "  → linux/arm64..."
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)
	@echo "  → windows/amd64..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	@echo "✓ All binaries built in $(BUILD_DIR)/"
	@ls -lh $(BUILD_DIR)/

.PHONY: install
install: build ## Install binary to /usr/local/bin
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✓ Installed: $(INSTALL_PATH)/$(BINARY_NAME)"

.PHONY: uninstall
uninstall: ## Remove binary from /usr/local/bin
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✓ Uninstalled"

# =============================================================================
# Run Targets
# =============================================================================

.PHONY: run
run: ## Run the CLI (shows help)
	@go run $(MAIN_PACKAGE)

.PHONY: generate
generate: ## Generate all test images to ./test-images/
	@go run $(MAIN_PACKAGE) generate --output $(OUTPUT_DIR)

.PHONY: generate-quick
generate-quick: ## Generate minimal test set (platform ratios, medium size, jpeg only)
	@go run $(MAIN_PACKAGE) generate --ratios platform --sizes medium --formats jpeg --output $(OUTPUT_DIR)

.PHONY: generate-minimal
generate-minimal: ## Generate using minimal config
	@go run $(MAIN_PACKAGE) generate --config configs/minimal.json --output $(OUTPUT_DIR)

.PHONY: list
list: ## List all available presets
	@go run $(MAIN_PACKAGE) list

.PHONY: version
version: ## Show version information
	@go run $(MAIN_PACKAGE) --version

# =============================================================================
# Development Targets
# =============================================================================

.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@echo "✓ Dependencies downloaded"

.PHONY: tidy
tidy: ## Tidy go.mod and go.sum
	@echo "Tidying modules..."
	@go mod tidy
	@echo "✓ Modules tidied"

.PHONY: fmt
fmt: ## Format all Go source files
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

.PHONY: vet
vet: ## Run go vet on all packages
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ No issues found"

.PHONY: lint
lint: ## Run golangci-lint (requires golangci-lint installed)
	@echo "Running linter..."
	@golangci-lint run ./...
	@echo "✓ Lint passed"

.PHONY: check
check: fmt vet ## Run fmt and vet checks
	@echo "✓ All checks passed"

# =============================================================================
# Test Targets
# =============================================================================

.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	@go test ./... -v
	@echo "✓ Tests passed"

.PHONY: test-short
test-short: ## Run tests in short mode
	@go test ./... -short

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

.PHONY: test-race
test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	@go test ./... -race
	@echo "✓ No race conditions detected"

# =============================================================================
# Clean Targets
# =============================================================================

.PHONY: clean
clean: ## Remove built binary and generated test images
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(OUTPUT_DIR)
	@rm -f coverage.out coverage.html
	@echo "✓ Cleaned"

.PHONY: clean-dist
clean-dist: ## Remove dist directory
	@rm -rf $(BUILD_DIR)

.PHONY: clean-all
clean-all: clean clean-dist ## Remove all generated files (binary, dist, test images)
	@echo "✓ All cleaned"

# =============================================================================
# Info Targets
# =============================================================================

.PHONY: info
info: ## Show project information
	@echo ""
	@echo "$(CYAN)Project Information$(RESET)"
	@echo "  Binary:      $(BINARY_NAME)"
	@echo "  Version:     $(VERSION)"
	@echo "  Go Version:  $(GO_VERSION)"
	@echo "  Output Dir:  $(OUTPUT_DIR)"
	@echo "  Build Dir:   $(BUILD_DIR)"
	@echo ""
	@echo "$(CYAN)Source Files$(RESET)"
	@find . -name "*.go" -not -path "./vendor/*" | wc -l | xargs echo "  Go files:   "
	@echo ""
	@echo "$(CYAN)Dependencies$(RESET)"
	@go list -m all | tail -n +2 | head -10
	@echo "  ..."

