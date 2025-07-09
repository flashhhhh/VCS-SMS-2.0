# How to write Dockerfile for mounting this local directory to the container so that the code can be edited without rebuilding the image
FROM golang:alpine

WORKDIR /app

CMD ["sh", "-c", "go mod tidy && go run cmd/server/rest/main.go"]