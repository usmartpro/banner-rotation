# Собираем в гошке
FROM golang:1.17.3 as build

RUN go install github.com/antelman107/net-wait-go@latest

WORKDIR "/opt/banner"

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . "/opt/banner/"
