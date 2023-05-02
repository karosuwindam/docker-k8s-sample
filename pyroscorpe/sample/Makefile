TAG = 0.0.1
DOCKER = docker
NAME = bookserver2:31000/tool/pyroscorpe-sample

BUILD = buildctl
BUILD_ADDR = tcp://buildkit.bookserver.home:1234 #arm64
BUILD_ADDR_ARM = tcp://buildkit-arm.bookserver.home:1235 #arm
BUILD_OPTION = "type=image,push=true,registry.insecure=true"

TEMPLATE = ./Dockerfile_tmp
TARGET = Dockerfile
TARGET_FILE = ./
GO_VERSION = 1.18.2
ARCH = armv6l
ARCH = arm64
OPT = "--privileged"

create:
	@echo "--- ${NAME} ${TAG} create ---"
	@echo "--- create Dockerfile ---"
	@cat ${TEMPLATE} | sed s/TAG/${TAG}/ | sed s/ARCH/${ARCH}/ | sed s/GO_VERSION/${GO_VERSION}/ > ${TARGET_FILE}/${TARGET}
build: create
	@echo "--- build Dockerfile --"
	${DOCKER} build -t ${NAME}:${TAG} -f ${TARGET_FILE}/${TARGET} ${TARGET_FILE}
build-kit: create
	@echo "--- buildkit build --"
	${BUILD} --addr ${BUILD_ADDR} build --output name=${NAME}:${TAG},${BUILD_OPTION} --frontend=dockerfile.v0 --local context=${TARGET_FILE}   --local dockerfile=${TARGET_FILE} --opt source=${TARGET}
run:
	${DOCKER} run --rm --name=booknewread ${NAME}:${TAG}