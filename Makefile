# Makefile for Online Canvas Games

# Variables
DOCKER_COMPOSE = docker compose
TSCOMPILER = ./tscompiler
TEMPL = templ

TS_GAMEASSETS = storageFiles/typescript/gameassets
TS_PAGES = storageFiles/typescript/pages

TEMPL_PATH = ./services/gateway/

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
	$(TSCOMPILER) $(TS_GAMEASSETS) $(TS_PAGES)

# Compile templ files
compile-templ:
	@echo "Compiling .templ files..."
	$(TEMPL) generate $(TEMPL_PATH)

# Clean up Docker resources
clean:
	@echo "Cleaning up Docker resources..."
	$(DOCKER_COMPOSE) down --rmi all --volumes --remove-orphans

.PHONY: build up down compile-ts compile-templ clean