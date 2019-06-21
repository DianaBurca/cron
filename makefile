BINARY = cron
GOARCH = amd64

VERSION?=?

CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})

DOCKER_IMAGE_NAME       ?= ${BINARY}
DOCKER_IMAGE_TAG        ?= latest
LDFLAGS = -ldflags "-w -s -X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"
all: linux docker
SHELL := /bin/bash

clean:
	go clean -n
	rm -f ${CURRENT_DIR}/${BINARY}-linux-${GOARCH}

linux:
	go get ./...
	@echo ">> building linux binary"
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} . ;

build_tag_push:
	@echo "Prepare linux build...";
	make linux;
	mv ${BINARY}-linux-${GOARCH} build/${BINARY}-linux-${GOARCH};
	@echo "Build docker image";
	docker build -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" -f build/Dockerfile ./build;
	@echo "Prepare image to be pushed";
	docker tag $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) diana1996/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG);
	@echo "Dockerhub login..."
	docker login -u=diana1996 -p=${DOCKERHUB_PASSWD}
	@echo "Push image to registry";
	docker push diana1996/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG);


.PHONY: release all linux
