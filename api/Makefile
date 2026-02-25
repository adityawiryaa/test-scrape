.PHONY: build test clean run-controller run-agent run-worker run-all tidy fix

build:
	CGO_ENABLED=1 go build -o bin/controller ./cmd/controller
	CGO_ENABLED=0 go build -o bin/agent ./cmd/agent
	CGO_ENABLED=0 go build -o bin/worker ./cmd/worker

test:
	go test ./... -v -race

clean:
	rm -rf bin/
	rm -f controller.db

tidy:
	go mod tidy

fix:
	go fix ./...

run-controller:
	go run ./cmd/controller

run-agent:
	go run ./cmd/agent

run-worker:
	go run ./cmd/worker

run-all:
	@echo "Starting controller on :6001..."
	@go run ./cmd/controller &
	@sleep 2
	@echo "Starting worker on :6002..."
	@go run ./cmd/worker &
	@sleep 2
	@echo "Starting agent..."
	@go run ./cmd/agent &
	@echo "All services started. Press Ctrl+C to stop."
	@wait

docker-controller:
	docker compose -f deployments/controller/docker-compose.yml up --build

docker-agent-worker:
	docker network create config-manager || true
	docker compose -f deployments/agent-worker/docker-compose.yml up --build

swagger:
	swag init -g cmd/controller/main.go -o docs/controller
	swag init -g cmd/worker/main.go -o docs/worker
