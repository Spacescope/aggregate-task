COMMIT_ID = $(shell git rev-parse --short HEAD)
ifeq ($(COMMIT_ID),)
COMMIT_ID = 'latest'
endif

.PHONY: test
IMAGE_PREFIX ?= 129862287110.dkr.ecr.us-east-2.amazonaws.com/data-infra
REGISTRY_SERVER ?= 129862287110.dkr.ecr.us-east-2.amazonaws.com/

help:
	@echo
	@echo "  binary - build binary"
	@echo "  build-aggregate-task - build docker images for centos"
	@echo "  swag - regenerate swag"
	@echo "  build-all - build docker images for centos"
	@echo "  push images to docker hub"

swag:
	swag init -g cmd/aggregate-task/main.go

binary:
	go build -o bin/aggregate-task cmd/aggregate-task/main.go

test:
	go clean -testcache
	gotestsum --format pkgname

build-aggregate-task:
	docker build -t $(IMAGE_PREFIX)/aggregate-task:$(COMMIT_ID) -f build/Dockerfile .

build-all: build-aggregate-task

push:
	docker push $(IMAGE_PREFIX)/aggregate-task:$(COMMIT_ID)
