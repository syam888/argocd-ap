SHELL := /bin/bash
BINARY := ap
PKG := ./cmd/ap
BUILD_DIR := bin

.PHONY: all build clean fmt

all: build

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) $(PKG)

clean:
	@rm -rf $(BUILD_DIR)

fmt:
	go fmt ./...
