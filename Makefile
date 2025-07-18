# Go parameters
GO ?= go
PACKAGE_NAME := daily_dice_roller
OUTPUT_DIR := ./build

all: build

# Build a static library (.a)
build-static:
	mkdir -p $(OUTPUT_DIR)
	$(GO) build -o . -o build

# Build both shared and static libraries
build: build-static

# Run tests
test:
	$(GO) test ./...

# Clean build artifacts
clean:
	rm -rf $(OUTPUT_DIR) *.so *.a *.h

.PHONY: all build build-shared build-static test clean
