ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run: modules up migrates
	go run cmd/orderservice/main.go

modules:
	go mod tidy

lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.37.1 golangci-lint run

fmt:
	go fmt ./...

test:
	go test ./...

up:
	docker-compose -f docker/docker-compose.yml up -d

down:
	docker-compose -f docker/docker-compose.yml down

migration:
	echo -n "Key: "; read MIGRATION_NAME; \
	docker run --rm -v $(shell pwd)/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq "$$MIGRATION_NAME"
	sudo chown $$USER:$$USER ./migrations/*

migrates:
	docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate \
        -path=/migrations -database mysql://$(DATABASE_USER):$(DATABASE_PASSWORD)@/$(DATABASE_NAME) up