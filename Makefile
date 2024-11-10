PKG := ./...

.PHONY: all test build run stop clean

all: docker-up

test:
	@echo "Running tests..."
	go test $(PKG)

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