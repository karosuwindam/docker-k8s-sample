# TAG = 24.0.0-fpm-alpine
TAG = 0.12b
DOCKER = docker
NAME = bookserver2:31000/moni/dockeri2cmoni

TARGET = ./Dockerfile
TARGET_FILE = ./services/app
GO_VERSION = 1.18.2
ARCH = ${shell docker version --format '{{.Server.Arch}}'}

BUILD = buildctl
BUILD_ADDR = tcp://buildkit.bookserver.home:1234 #arm64
BUILD_ADDR_ARM = tcp://buildkit-arm.bookserver.home:1235 #arm
BUILD_OPTION = "type=image,push=true,registry.insecure=true"
GO_VERSION = 1.18.2
DOCKERFILE_FOLDER = .

all: build push
create:
	@echo "--- create Dockerfile ---"
	@cat Dockerfile_tmp | sed s/TAG/${TAG}/ | sed s/ARCH/${ARCH}/ | sed s/GO_VERSION/${GO_VERSION}/ > ${TARGET}
build: create
	@echo "--- build Dockerfile ---"
	${DOCKER} build -t ${NAME}:${TAG} -f ${TARGET} ${TARGET_FILE}
buildkit: create
	@echo "--- buildkit Dockerfile ---"
	${BUILD} --addr ${BUILD_ADDR_ARM} build --output name=${NAME}:${TAG},${BUILD_OPTION} --frontend=dockerfile.v0 --local context=${TARGET_FILE}   --local dockerfile=${DOCKERFILE_FOLDER}
push:
	${DOCKER} push ${NAME}:${TAG}
help:
	@echo "make |all|create|build|push"