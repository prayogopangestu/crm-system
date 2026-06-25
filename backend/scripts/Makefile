.PHONY: run build test race lint migrate-up migrate-down migrate-status seed proto compose-up compose-down

run:
	go run ./cmd/api

build:
	go build ./...

test:
	go test ./...

race:
	go test -race ./...

lint:
	go vet ./...
	golangci-lint run

migrate-up:
	go run ./cmd/migrate -direction up

migrate-down:
	go run ./cmd/migrate -direction down

migrate-status:
	go run ./cmd/migrate -direction status

seed:
	go run ./cmd/seed

proto:
	docker run --rm -v "$(CURDIR)/..:/workspace" -w /workspace bufbuild/buf:1.47.2 generate

compose-up:
	docker compose -f deployments/docker-compose.yml up --build

compose-down:
	docker compose -f deployments/docker-compose.yml down

