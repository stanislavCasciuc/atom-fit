FROM golang:1.23.1

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download


RUN go install github.com/swaggo/swag/cmd/swag@latest 

COPY . .
RUN swag init -g ./main/main.go -d cmd,api,internal && swag fmt



RUN go build -o /app/bin/social cmd/main/main.go

EXPOSE 8080

CMD ["/app/bin/social"]
