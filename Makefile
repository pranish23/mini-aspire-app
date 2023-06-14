# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean

# Build target
BINARY_NAME=miniAspireApp

# Build the Go program
build:
	$(GOBUILD) -o $(BINARY_NAME) 

# Run the Go program using 'go run main.go'
dev:
	$(GOCMD) run main.go

# Run the Go binary
run:
	./$(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./test/...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f ./$(BINARY_NAME)

.PHONY: build dev run test clean
