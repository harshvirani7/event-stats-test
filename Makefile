# Makefile for building and running the Go service

.PHONY: build run stop clean

SERVICE_NAME := event-stats-test

build:
    @echo "Building $(SERVICE_NAME)..."
    go build -o $(SERVICE_NAME) .

run:
    @echo "Starting $(SERVICE_NAME)..."
    ./$(SERVICE_NAME)

stop:
    @echo "Stopping $(SERVICE_NAME)..."
    # If your service needs to be stopped in a specific way, add the command here

clean:
    @echo "Cleaning up..."
    go clean
    rm -f $(SERVICE_NAME)
