BIN=bin/rotator.app
DOCKER_IMG=rotator
RELEASE=develop

PACKAGE_PATH := "github.com/FedoseevAlex/banner-rotation"
GIT_HASH := $(shell git rev-parse HEAD)
LDFLAGS := -X $(PACKAGE_PATH)/internal/common.release="develop"
LDFLAGS += -X $(PACKAGE_PATH)/internal/common.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S)
LDFLAGS += -X $(PACKAGE_PATH)/internal/common.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/...

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG):$(RELEASE) \
		-f build/Dockerfile .

run-img: build-img
	docker rm $(DOCKER_IMG) || true
	docker run \
	--name $(DOCKER_IMG) \
	-p 8080:8080 \
	$(DOCKER_IMG):$(RELEASE)

version: build
	$(BIN) version

test:
	go test -race -count 100 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
