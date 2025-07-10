FROM golang:alpine AS builder

RUN apk add --no-cache ca-certificates && \
    update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o rest-server cmd/server/rest/main.go && \
    go build -o grpc-server cmd/server/grpc/main.go && \
    go build -o kafka-server cmd/server/kafka/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates && \
    update-ca-certificates

WORKDIR /app

COPY --from=builder /app/rest-server ./rest-server
COPY --from=builder /app/grpc-server ./grpc-server
COPY --from=builder /app/kafka-server ./kafka-server
COPY configs/ ./configs/

CMD ["sh", "-c", "./rest-server & ./grpc-server & ./kafka-server"]