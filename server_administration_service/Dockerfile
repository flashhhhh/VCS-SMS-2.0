FROM golang:alpine

WORKDIR /app

CMD ["sh", "-c", "go mod tidy && go run cmd/server/rest/main.go & go run cmd/server/grpc/main.go"]