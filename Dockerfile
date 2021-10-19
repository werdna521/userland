FROM golang:alpine

WORKDIR /app/api

COPY . .

RUN go mod download

ENTRYPOINT ["go", "run", "main.go"]
