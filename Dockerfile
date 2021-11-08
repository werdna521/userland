FROM golang:alpine

RUN apk add alpine-sdk
RUN apk add build-base

WORKDIR /app/api
COPY . .

RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -a -v -tags musl -o /userland

# RUN go install github.com/mitranim/gow@latest
# ENTRYPOINT ["/go/bin/gow", "run", "main.go", '-a', '-v', '-tags', 'musl']
ENTRYPOINT ["/userland"]
