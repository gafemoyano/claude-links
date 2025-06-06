# Simple Makefile for a Go project

# Build the application
all: swag build

# swag:
# 	@echo "Generating Swagger docs..."
# 	@if [ -z "$(shell find . -name "*.go" -newer ./docs/swagger.json 2>/dev/null)" ] && [ -f ./docs/swagger.json ]; then \
# 		echo "Swagger docs are up to date"; \
# 	else \
# 		swag init; \
# 	fi


build:
	@echo "Building..."
	@go build -o main main.go

# Run the application
run:
	@go run main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air;\
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean
