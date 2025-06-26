FROM golang:alpine

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod tidy

CMD ["sh", "-c", "go run cmd/server/rest/main.go & go run cmd/server/grpc/main.go"]