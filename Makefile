build:
	@go build -o bin/social cmd/main/main.go

run: build
	@./bin/social
