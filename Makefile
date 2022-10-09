GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
BIN := "./bin/top"
CLIENT_BIN := "./bin/client"

# Linter
install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.48.0

lint: install-lint-deps
	golangci-lint run ./...

test:
	go test ./... -race -count 100

# All
build: build-top build-client

# Top
build-top:
	CGO_ENABLED=0 GOOS=linux go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/top

run-top: build-top
	$(BIN) -config ./configs/app.toml --grpc-addr=:9990

# Client
build-client:
	CGO_ENABLED=0 GOOS=linux go build -v -o $(CLIENT_BIN) -ldflags "$(LDFLAGS)" ./cmd/client
run-client: build-client
	$(CLIENT_BIN) -config ./configs/client.toml -n 5 -m 15 --grpc-addr=:9990

# Dev
up-dev: build
	docker-compose -f ./deployments/docker-compose.dev.yaml up -d --build
down-dev:
	docker-compose -f ./deployments/docker-compose.dev.yaml down --remove-orphans
logs-dev:
	docker-compose -f deployments/docker-compose.dev.yaml logs -f
rest-dev: down-dev up-dev logs-dev

# GRPC
install-protobuf:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
generate:
	protoc ./api/TopService.proto --go_out=./api --go-grpc_out=./api

.PHONY: run build
