FROM golang:alpine AS builder

RUN apk add --no-cache ca-certificates && \
    update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o rest-server cmd/server/rest/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates && \
    update-ca-certificates

WORKDIR /app

COPY --from=builder /app/rest-server ./rest-server
COPY configs/ ./configs/

CMD ["./rest-server"]