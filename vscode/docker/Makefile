TAG = 1.75.1
BUILD = buildctl
BUILD_ADDR = tcp://buildkit.bookserver.home:1234 #arm64
BUILD_ADDR_ARM = tcp://buildkit-arm.bookserver.home:1235 #arm
NAME = bookserver2:31000/karosu/vscode-go
BUILD_OPTION = "type=image,push=true,registry.insecure=true"
DOCKERFILE_FOLDER = .
GO_VERSION = 1.19.6
TARGET = ${DOCKERFILE_FOLDER}/Dockerfile
ARCH = arm64
GO_FILE = "go${GO_VERSION}.linux-${ARCH}.tar.gz"
# ARCH = armv6l
all: build
create:
	@echo "--- create Dockerfile ---"
	@echo "cat Dockerfile_tmp | sed s/"%TAG%"/${TAG}/ | sed s/"%GO_VERSION%"/${GO_VERSION}/ | sed s/"%GO_FILE%"/${GO_FILE}/  > ${TARGET}"
	@cat Dockerfile_tmp | sed s/"%TAG%"/${TAG}/ | sed s/"%GO_VERSION%"/${GO_VERSION}/ | sed s/"%GO_FILE%"/${GO_FILE}/  > ${TARGET}
build: create
	@echo "--- build Dockerfile --"
	${BUILD} --addr ${BUILD_ADDR} build --output name=${NAME}:${TAG},${BUILD_OPTION} --frontend=dockerfile.v0 --local context=${DOCKERFILE_FOLDER}   --local dockerfile=${DOCKERFILE_FOLDER}
