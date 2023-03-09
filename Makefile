all: migrate test vet

migrate:
	go run ./cmd/migrate/main.go

test:
	go test ./...

vet:
	go vet ./...

MOCKS_DESTINATION=mocks
.PHONY: mocks
mocks: database/drugs.go database/illnesses.go database/procedures.go database/scheme_days.go database/schemes.go database/treatments.go database/treatment_schemes.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done
