FROM golang:alpine

WORKDIR /app

CMD ["sh", "-c", "go mod tidy && go run cmd/main.go"]