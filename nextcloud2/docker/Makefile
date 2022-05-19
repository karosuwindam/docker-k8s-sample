# TAG = 24.0.0-fpm-alpine
TAG = 24.0.0
DOCKER = docker
NAME = bookserver2:31000/karosu/nextcloud

TARGET = ./Dockerfile
ARCH = ${shell docker version --format '{{.Server.Arch}}'}

all: build push
build:
	@echo "--- create Dockerfile ---"
	@cat Dockerfile_tmp | sed s/TAG/${TAG}/ | sed s/ARCH/${ARCH}/ > ${TARGET}
	@echo "--- build Dockerfile --"
	${DOCKER} build -t ${NAME}:${TAG} -f ${TARGET} .
push:
	${DOCKER} push ${NAME}:${TAG}
help:
	@echo "make |all|build|push"