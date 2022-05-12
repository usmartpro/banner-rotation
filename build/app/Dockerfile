# Собираем в гошке
FROM golang:1.17.3 as build

ENV BIN_FILE /opt/banner/app
ENV CODE_DIR /go/src/
ENV CGO_ENABLED=0

RUN go install github.com/antelman107/net-wait-go@latest

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/banner/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL SERVICE="banner-rotation"
LABEL MAINTAINERS="usmartpro@gmail.com"

ENV BIN_FILE "/opt/banner/app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

COPY --from=build /go/bin/net-wait-go /bin/net-wait-go

ENV ENV_FILE /go/src/.env
COPY --from=build ./.env ${ENV_FILE}

CMD ${BIN_FILE} -env ${ENV_FILE}
