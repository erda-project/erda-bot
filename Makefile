DOCKER_REGISTRY ?= registry.erda.cloud/erda
BUILD_TIME := $(shell date "+%Y%m%d-%H%M%S")
COMMIT_ID := $(shell git rev-parse --short HEAD 2>/dev/null)
IMG ?= ${DOCKER_REGISTRY}/erda-bot:${BUILD_TIME}-${COMMIT_ID}

build-on-local:
	docker build --platform=linux/amd64 -t ${IMG} . --push --build-arg GOPROXY=""
	echo 'image: ${IMG}'

lint:
	env GOGC=25 golangci-lint run --fix -j 8 -v ./...

