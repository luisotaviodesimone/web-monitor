FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go install github.com/jackc/tern/v2@latest
