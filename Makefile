# Go parameters
GO ?= go
PACKAGE_NAME := daily_dice_roller
OUTPUT_DIR := ./build

all: build

build-static:
	mkdir -p $(OUTPUT_DIR)
	$(GO) build -o . -o build

# Build both shared and static libraries
build: build-static

# Run tests
test:
	$(GO) test ./... -coverprofile=coverage.out

run: 
	$(GO) run main.go

coverage:
	$(GO) tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -rf $(OUTPUT_DIR) *.so *.a *.h coverage.out

.PHONY: all build build-shared build-static test clean coverage
