DOCKER_COMPOSE := docker-compose.yml
PKG := ./...

.PHONY: all test clean docker-up docker-down

all: docker-up

# Run tests
test:
	@echo "Running tests..."
	go test $(PKG)