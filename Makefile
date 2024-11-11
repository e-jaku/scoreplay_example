PKG := ./...
COVERAGE_FILE := coverage.out

.PHONY: all test test-coverage build run stop clean

all: docker-up

test:
	@echo "Running tests..."
	go test $(PKG)

test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=$(COVERAGE_FILE) $(PKG)	

build:
	@echo "Running build..."
	go build -o ./bin/scoreplay ./cmd/scoreplay/main.go

run: build
	@echo "Running application..."
	docker compose up -d && ./bin/scoreplay

stop:
	@echo "Running stop..."
	docker compose down

clean:
	@echo "Running clean..."
	@rm -rf bin