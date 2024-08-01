# Makefile for Online Canvas Games

# Variables
DOCKER_COMPOSE = docker compose
TSCOMPILER = ./tscompiler

# Default target
all: build up

# Build the Docker containers
build:
	@echo "Building Docker containers..."
	$(DOCKER_COMPOSE) build

# Start the Docker containers
up:
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up

# Stop the Docker containers
down:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

# Compile TypeScript files
compile-ts:
	@echo "Compiling TypeScript files..."
	$(TSCOMPILER) storageFiles/typescript/gameassets storageFiles/typescript/pages

# Clean up Docker resources
clean:
	@echo "Cleaning up Docker resources..."
	$(DOCKER_COMPOSE) down --rmi all --volumes --remove-orphans

.PHONY: all build up down compile-ts clean