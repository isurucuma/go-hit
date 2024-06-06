DEFAULT_GOAL := build

.PHONY: build run

APP_NAME=hit

build: 
	@echo "Building..."
	go build -o bin/$(APP_NAME) cmd/$(APP_NAME)/$(APP_NAME).go

run:
	@echo "Running..."
	go run cmd/$(APP_NAME)/$(APP_NAME).go