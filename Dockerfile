FROM golang:1.23.2 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go generate ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

FROM debian:12-slim

WORKDIR /app

RUN apt-get -y update && \
    apt-get -y upgrade && \
    apt-get -y install ca-certificates iputils-ping net-tools netcat-traditional procps tzdata && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/server /app/server
COPY cmd/api/.env /app

EXPOSE 8080

ENTRYPOINT ["/app/server"]