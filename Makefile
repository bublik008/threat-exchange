.PHONY: generate test lint run docker-up docker-down

generate:
		go generate ./...

test:
		go test ./...

lint:
		go vet ./...
		golangci-lint run

run:
		go run ./cmd/server

docker-up:
		docker compose -f deployments/docker-compose.yml up --build

docker-down:
		docker compose -f deployments/docker-compose.yml down
