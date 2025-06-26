FROM golang:alpine AS builder

RUN apk add --no-cache ca-certificates && \
    update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates && \
    update-ca-certificates

WORKDIR /app

COPY --from=builder /app/main ./main
COPY configs/ ./configs/

CMD ["./main"]