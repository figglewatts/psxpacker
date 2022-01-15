GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
BINARY_NAME=psxpacker

all: build
build: 
	$(GOBUILD) -o bin/$(BINARY_NAME) -v cmd/$(BINARY_NAME)/main.go
install:
	$(GOINSTALL)
clean: 
	$(GOCLEAN)
	rm -rf bin