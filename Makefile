GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
BIN := "./bin/top"

# Linter
install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.48.0

lint: install-lint-deps
	golangci-lint run ./...

test:
	go test ./...

# Top
build:
	CGO_ENABLED=0 GOOS=linux go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/top

run: build
	$(BIN) -config ./configs/app.toml

# Dev
up-dev: build
	docker-compose -f ./deployments/docker-compose.dev.yaml up -d --build
down-dev:
	docker-compose -f ./deployments/docker-compose.dev.yaml down --remove-orphans
dev-run:
	docker-compose -f ./deployments/docker-compose.dev.yaml run ubunty-testing bash

# GRPC
install-protobuf:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
generate:
	protoc ./api/TopService.proto --go_out=./api --go-grpc_out=./api

.PHONY: run build