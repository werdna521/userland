FROM golang:alpine

WORKDIR /app/api

COPY . .

RUN go mod download

# install golang-migrate
RUN curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
RUN echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ RUN(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
RUN apt-get update
RUN apt-get install -y migrate

RUN migrate -database \"postgres://postgres:password@postgres/userland?sslmode=disable\" -path db/migrations up

ENTRYPOINT ["go", "run", "main.go"]
