TAG = 0.14.5
DOCKER = docker
NAME = bookserver2:31000/moni/dockeri2cmoni
LABEL_NAME = "senser eijent"

BUILD = buildctl
BUILD_ADDR = tcp://buildkit.bookserver.home:1234 #arm64
BUILD_ADDR_ARM = tcp://buildkit-arm.bookserver.home:1235 #arm
BUILD_OPTION = "type=image,push=true,registry.insecure=true"

TARGET = ./Dockerfile
TEMPLATE = Dockerfile_tmp
TARGET_FILE = ./golang
GO_VERSION = 1.22.1
ARCH = arm64
DOCKERFILE_FOLDER = .

all: build push
create:
	@echo "--- ${NAME} ${TAG} create ---"
	@echo "--- create Dockerfile ---"
	@cat ${TEMPLATE} | sed s/TAG/${TAG}/ | sed s/ARCH/${ARCH}/ | sed s/GO_VERSION/${GO_VERSION}/ > ${TARGET}
build: create
	@echo "--- build Dockerfile --"
	${DOCKER} build -t ${NAME}:${TAG} -f ${TARGET} ${TARGET_FILE}
push:
	${DOCKER} push ${NAME}:${TAG}
build-kit: create
	@echo "--- buildkit build --"
	${BUILD} --addr ${BUILD_ADDR} build --output name=${NAME}:${TAG},${BUILD_OPTION} --frontend=dockerfile.v0 --local context=${TARGET_FILE}   --local dockerfile=${DOCKERFILE_FOLDER}
build-kit-arm: create
	@echo "--- buildkit build --"
	${BUILD} --addr ${BUILD_ADDR_ARM} build --output name=${NAME}:${TAG},${BUILD_OPTION} --frontend=dockerfile.v0 --local context=${TARGET_FILE}   --local dockerfile=${DOCKERFILE_FOLDER}
help:
	@echo "make |all|create|build|push"
