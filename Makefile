all: migrate test vet

migrate:
	go run ./cmd/migrate/main.go

test:
	go test ./...

vet:
	go vet ./...

MOCKS_DESTINATION=mocks
.PHONY: mocks
# put the files with interfaces you'd like to mock in prerequisites
# wildcards are allowed
mocks: database/drugs.go database/illnesses.go database/procedures.go database/scheme_days.go database/schemes.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done
