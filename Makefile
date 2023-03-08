all: migrate test vet

migrate:
	go run ./cmd/migrate/main.go

test:
	go test ./...

vet:
	go vet ./...
