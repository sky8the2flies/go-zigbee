# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=go-zigbee-herdsman

# Directories
SRC_DIR=./cmd
BIN_DIR=./bin

# Targets
all: test build

build:
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) -v $(SRC_DIR)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BIN_DIR)/$(BINARY_NAME)

deps:
	$(GOGET) -v ./...

.PHONY: all build test clean deps