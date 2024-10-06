include .env
MIGRATION_PATH = ./db/migrations
gen-docs:
	@swag init -g ./main/main.go -d cmd,api,internal && swag fmt

build:
	@go build -o bin/social cmd/main/main.go

run: gen-docs build
	@./bin/social

migration:
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(filter-out $@,$(MAKECMDGOALS))
migrate-up:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) up	
migrate-down:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))	$(filter-out $@,$(MAKECMDGOALS))

