TAG = 0.0.2
DOCKER = docker
NAME = bookserver2:31000/tool/gotext-server
APPNAME = "gotext-server"
GOVERSION = "1.18"
TEMPLATE = ./Dockerfile_tmp
BASE_CONTANER = "debian:11"

BUILD = buildctl
BUILD_ADDR = tcp://buildkit.bookserver.home:1234 #arm64
BUILD_ADDR_ARM = tcp://buildkit-arm.bookserver.home:1235 #arm
BUILD_OPTION = "type=image,push=true,registry.insecure=true"

TEMPLATE = ./gotext_server/Dockerfile_tmp
TARGET = Dockerfile
TARGET_FILE = ./gotext_server
GO_VERSION = 1.18.2
ARCH = armv6l
OPT = "--privileged"

all: build push
create:
	@echo "--- ${NAME} ${TAG} create ---"
	@echo "--- create Dockerfile ---"
	@cat ${TEMPLATE} | sed -e "s|TAG|${TAG}|g" | sed -e "s|GOVERSION|${GOVERSION}|g" | sed -e "s|APPNAME|${APPNAME}|g" | sed -e "s|BASE_CONTANER|${BASE_CONTANER}|g" > ${TARGET_FILE}/${TARGET}
build: create
	@echo "--- build Dockerfile --"
	${DOCKER} build -t ${NAME}:${TAG} -f ${TARGET_FILE}/${TARGET} ${TARGET_FILE}
push:
	${DOCKER} push ${NAME}:${TAG}
rmi:
	${DOCKER} rmi ${NAME}:${TAG}
	${DOCKER} image prune -f
build-kit: create
	@echo "--- buildkit build --"
	${BUILD} --addr ${BUILD_ADDR_ARM} build --output name=${NAME}:${TAG},${BUILD_OPTION} --frontend=dockerfile.v0 --local context=${TARGET_FILE}   --local dockerfile=${TARGET_FILE} --opt source=${TARGET_FILE}/${TARGET}
run:
	${DOCKER} run --rm --name=gotext-server -p 18080:8080 ${NAME}:${TAG}
help:
	@echo "make |all|create|build|push"