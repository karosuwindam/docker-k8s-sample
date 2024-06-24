# TAG = 24.0.0-fpm-alpine
TAG = 0.22.0
DOCKER = docker
NAME = bookserver2:31000/karosu/booknewread

BUILD = buildctl
BUILD_ADDR = tcp://buildkit.bookserver.home:1234 #arm64
BUILD_ADDR_ARM = tcp://buildkit-arm.bookserver.home:1235 #arm
BUILD_OPTION = "type=image,push=true,registry.insecure=true"

TEMPLATE = ./app/Dockerfile_tmp
TARGET = Dockerfile_arm
TARGET_FILE = ./app
GO_VERSION = 1.22.1
ARCH = aarch64
OPT = "--privileged"

all: build push
create:
	@echo "--- ${NAME} ${TAG} create ---"
	@echo "--- create Dockerfile ---"
	@cat ${TEMPLATE} | sed s/TAG/${TAG}/ | sed s/ARCH/${ARCH}/ | sed s/GO_VERSION/${GO_VERSION}/ > ${TARGET_FILE}/${TARGET}
build: create
	@echo "--- build Dockerfile --"
	${DOCKER} build -t ${NAME}:${TAG} -f ${TARGET_FILE}/${TARGET} ${TARGET_FILE}
push:
	${DOCKER} push ${NAME}:${TAG}
build-kit: create
	@echo "--- buildkit build --"
	${BUILD} --addr ${BUILD_ADDR} build --output name=${NAME}:${TAG},${BUILD_OPTION} --frontend=dockerfile.v0 --local context=${TARGET_FILE}   --local dockerfile=${TARGET_FILE} --opt source=${TARGET}
run:
	${DOCKER} run --rm --name=booknewread -p 18080:8080 ${NAME}:${TAG}
help:
	@echo "make |all|create|build|push"