# TAG = 24.0.0-fpm-alpine
TAG = 0.13.0
DOCKER = docker
NAME = bookserver2:31000/karosu/booknewread

TEMPLATE = ./app/Dockerfile_tmp
TARGET = ./app/Dockerfile_arm
TARGET_FILE = ./app
GO_VERSION = 1.18.2
ARCH = ${shell docker version --format '{{.Server.Arch}}'}

all: build push
create:
	@echo "--- create Dockerfile ---"
	@cat ${TEMPLATE} | sed s/TAG/${TAG}/ | sed s/ARCH/${ARCH}/ | sed s/GO_VERSION/${GO_VERSION}/ > ${TARGET}
build: create
	@echo "--- build Dockerfile --"
	${DOCKER} build -t ${NAME}:${TAG} -f ${TARGET} ${TARGET_FILE}
push:
	${DOCKER} push ${NAME}:${TAG}
help:
	@echo "make |all|create|build|push"