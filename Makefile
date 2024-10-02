include .env
MIGRATION_PATH = ./db/migrations

build:
	@go build -o bin/social cmd/main/main.go

run: build
	@./bin/social

migration:
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(filter-out $@,$(MAKECMDGOALS))
migrate-up:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) up	
migrate-down:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))	$(filter-out $@,$(MAKECMDGOALS))

