# Go parameters
GO=go
SERVICES=attendanceops salaryops
PKGS=$(shell go list ./... | grep -v /cmd/)  # All packages, excluding services in cmd/

# Build commands
.PHONY: all build clean run

all: build

build: $(SERVICES)
$(SERVICES):
	@echo "Building $@..."
	$(GO) build -o bin/$@ ./cmd/$@

clean:
	@echo "Cleaning up binaries..."
	rm -rf bin/*

run: $(SERVICES)
run-%: bin/%
	@echo "Running $*..."
	./bin/$*

# Test commands
.PHONY: test test-all

test:
	@echo "Running tests..."
	$(GO) test -v ./...

test-%:
	@echo "Running tests for $*..."
	$(GO) test -v ./cmd/$*

test-all:
	@echo "Running tests for all packages..."
	$(GO) test -v $(PKGS)

# Linting (requires golangci-lint to be installed)
.PHONY: lint lint-all

lint:
	@echo "Linting all packages..."
	golangci-lint run $(PKGS)

lint-%:
	@echo "Linting $*..."
	golangci-lint run ./cmd/$*

# Utility targets
.PHONY: deps

deps:
	@echo "Installing dependencies..."
	$(GO) mod download
