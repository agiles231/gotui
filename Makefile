# Go project Makefile

# Variables
BINARY_NAME := gotui
MODULE := github.com/agiles231/gotui
GO := go
GOFLAGS := -v
LDFLAGS := -s -w

# Build directory
BUILD_DIR := build

.PHONY: all build run test clean fmt lint vet tidy help

## all: Build the binary (default target)
all: build

## build: Build the binary
build:
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .

## run: Run the application
run:
	$(GO) run .

## test: Run all tests
test:
	$(GO) test $(GOFLAGS) ./...

## test-cover: Run tests with coverage
test-cover:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

## clean: Remove build artifacts
clean:
	$(GO) clean
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## fmt: Format code
fmt:
	$(GO) fmt ./...

## lint: Run golangci-lint (requires golangci-lint installed)
lint:
	golangci-lint run ./...

## vet: Run go vet
vet:
	$(GO) vet ./...

## tidy: Tidy and verify module dependencies
tidy:
	$(GO) mod tidy
	$(GO) mod verify

## install: Install the binary
install:
	$(GO) install .

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/ /'

