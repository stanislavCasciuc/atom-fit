FROM golang:1.23.1

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go install github.com/swaggo/swag/cmd/swag@latest 

COPY . .
RUN swag init -g ./main/main.go -d cmd,api,internal && swag fmt

ENV MIGRATION_PATH ./db/migrations

RUN /bin/sh -c "migrate -path=$MIGRATION_PATH -database=$DB_ADDR up"

RUN go build -o /app/bin/social cmd/main/main.go

EXPOSE 8080

CMD ["/app/bin/social"]
