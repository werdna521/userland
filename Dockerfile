FROM golang:alpine

WORKDIR /app/api

COPY . .

RUN go mod download

RUN go install github.com/mitranim/gow@latest

ENTRYPOINT ["/go/bin/gow", "run", "main.go"]
