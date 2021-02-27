all: modules check test run

run: modules
	go run main.go

modules:
	go mod tidy

check:
	golangci-lint run

fmt:
	go fmt ./...

test:
	go test ./...