# TAG = 24.0.0-fpm-alpine
TAG = 26.0.1
DOCKER = docker
NAME = bookserver2:31000/karosu/nextcloud

TARGET = ./Dockerfile
ARCH = ${shell docker version --format '{{.Server.Arch}}'}

all: build push
create:
	@echo "--- create Dockerfile ---"
	@cat Dockerfile_tmp | sed s/"%TAG%"/${TAG}/ | sed s/"%ARCH%"/${ARCH}/ > ${TARGET}
build: create
	@echo "--- build Dockerfile --"
	${DOCKER} build -t ${NAME}:${TAG} -f ${TARGET} .
push:
	${DOCKER} push ${NAME}:${TAG}
help:
	@echo "make |all|build|push"