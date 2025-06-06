# Simple Makefile for a Go project

# Build the application
all: swag build

swag:
	@echo "Generating Swagger docs..."
	@if [ -z "$(shell find . -name "*.go" -newer ./docs/swagger.json 2>/dev/null)" ] && [ -f ./docs/swagger.json ]; then \
		echo "Swagger docs are up to date"; \
	else \
		swag init; \
	fi



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

# Docker and deployment variables
PROJECT_ID ?= trii-dev-325214
IMAGE_NAME ?= be-links
REGION ?= us-east1
PORT ?= 3000

# Docker build for production
docker-build:
	@echo "Building Docker image..."
	@docker build -t gcr.io/$(PROJECT_ID)/$(IMAGE_NAME):latest .

# Push Docker image to Google Container Registry
docker-push:
	@echo "Pushing Docker image to GCR..."
	@docker push gcr.io/$(PROJECT_ID)/$(IMAGE_NAME):latest

# Deploy to Cloud Run
deploy:
	@echo "Deploying to Cloud Run..."
	@gcloud run deploy $(IMAGE_NAME) \
		--image gcr.io/$(PROJECT_ID)/$(IMAGE_NAME):latest \
		--region $(REGION) \
		--platform managed \
		--port $(PORT) \
		--allow-unauthenticated

# Build and push image (convenience target)
image-build-push: docker-build docker-push

# Full deployment pipeline
deploy-full: docker-build docker-push deploy

# Configure Docker to use gcloud as credential helper
docker-configure:
	@echo "Configuring Docker for GCR..."
	@gcloud auth configure-docker

.PHONY: all build run test clean docker-build docker-push deploy image-build-push deploy-full docker-configure
