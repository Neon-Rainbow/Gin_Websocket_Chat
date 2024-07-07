# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=myapp
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WIN=$(BINARY_NAME).exe
BINARY_MAC=$(BINARY_NAME)_mac

# Default target executed when no arguments are given to make.
default: build

# Build the project for different platforms
build: 
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_UNIX) -v
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_WIN) -v
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o bin/$(BINARY_MAC) -v

# Run the tests
test: 
	$(GOTEST) -v ./...

# Clean the directory
clean: 
	$(GOCLEAN)
	rm -f bin/$(BINARY_UNIX)
	rm -f bin/$(BINARY_WIN)
	rm -f bin/$(BINARY_MAC)

# Install the dependencies
deps:
	$(GOGET) -u ./...

# Format the code
fmt:
	$(GOCMD) fmt ./...

# Check for linting issues
lint:
	$(GOCMD) vet ./...

.PHONY: build test clean deps fmt lint