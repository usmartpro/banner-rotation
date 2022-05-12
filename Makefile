BIN := "./bin/banner"
DOCKER_IMG="banner-rotation:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

DB_CONN := "postgresql://postgres:postgres@localhost:5432/bnners?sslmode=disable"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/banner

run: build up
	#$(BIN)

version: build
	$(BIN) version

test:
	go test -race ./internal/... -count 100

test-int: build up
	set -e ;\
	test_status_code=0 ;\
	docker-compose run tests go test github.com/usmartpro/banner-rotation/cmd/tests/integration || test_status_code=$$? ;\
	docker-compose down ;\
	exit $$test_status_code ;

lint:
	golangci-lint run ./...

migrate:
	goose --dir=migrations postgres ${DB_CONN} up

up:
	docker-compose up -d

down:
	docker-compose down

generate:
	go generate ./...

.PHONY: build run build-img run-img version test lint
