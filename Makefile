.PHONY: build test clean deps vendor check-deps update-deps verify-deps clean-deps run dev

# Variables
BINARY_NAME=wctest
BUILD_DIR=build
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-s -w"

# Default target
all: build

# Create build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Download dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

# Vendor dependencies
vendor: deps
	$(GO) mod vendor

# Check for outdated dependencies
check-deps: deps
	$(GO) list -u -m all

# Update all dependencies
update-deps: deps
	$(GO) get -u ./...

# Verify dependencies
verify-deps: deps
	$(GO) mod verify

# Clean dependencies
clean-deps:
	rm -rf vendor/
	$(GO) clean -modcache

# Build the application
build: deps $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# Run the application
run: build
	$(BUILD_DIR)/$(BINARY_NAME)

# Development mode: rebuild and run
dev: clean build
	$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test: deps
	$(GO) test $(GOFLAGS) ./...

# Run tests with coverage
test-coverage: deps
	$(GO) test $(GOFLAGS) -coverprofile=$(BUILD_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(BUILD_DIR)/coverage.out

# Clean build artifacts
clean: clean-deps
	rm -rf $(BUILD_DIR)
	rm -f *.db 