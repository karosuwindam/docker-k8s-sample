TAG = 22.04
DOCKER = docker
ifndef ARCH
$(error ARCH is not set arm or arm64)
endif
ifeq (${ARCH},arm)
BUILD_ADDR = tcp://buildkit-arm.bookserver.home:1235 #arm
else ifeq (${ARCH},arm64)
BUILD_ADDR = tcp://buildkit.bookserver.home:1234 #arm64
else
$(error ARCH is not set arm or arm64)
endif
NAME = bookserver2:31000/${ARCH}/ubuntu

TEMPLATE = ./Dockerfile_tmp
TARGET = Dockerfile
TARGET_FILE = ./
GO_VERSION = 1.18.2


BUILD = buildctl
BUILD_OPTION = "type=image,push=true,registry.insecure=true"



ARCH = ${shell arch}
ifeq (${ARCH},x86_64)
ARCH = amd64
else
ARCH = armv6l
endif

OPT = "--privileged"

create:
	@echo "create dockerfile"
	@echo "for ${NAME}:${TAG}"
	@cat ${TEMPLATE} | sed s/TAG/${TAG}/ | sed s/ARCH/${ARCH}/ | sed s/GO_VERSION/${GO_VERSION}/ > ${TARGET_FILE}${TARGET}
build: create
	@echo "--- buildkit build --"
	${BUILD} --addr ${BUILD_ADDR} build --output name=${NAME}:${TAG},${BUILD_OPTION} --frontend=dockerfile.v0 --local context=${TARGET_FILE}   --local dockerfile=${TARGET_FILE} --opt source=${TARGET_FILE}${TARGET}
