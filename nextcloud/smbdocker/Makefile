TAG = 26.0.2
BUILD = buildctl
BUILD_ADDR = tcp://buildkit.bookserver.home:1234 #arm64
BUILD_ADDR_ARM = tcp://buildkit-arm.bookserver.home:1235 #arm
NAME = bookserver2:31000/karosu/nextcloud
BUILD_OPTION = "type=image,push=true,registry.insecure=true"
DOCKERFILE_FOLDER = .
TARGET = ${DOCKERFILE_FOLDER}/Dockerfile

all: build
create:
	@echo "--- create Dockerfile ---"
	@cat Dockerfile_tmp | sed s/"%TAG%"/${TAG}/ | sed s/"%ARCH%"/${ARCH}/ > ${TARGET}
build: create
	@echo "--- build Dockerfile --"
	${BUILD} --addr ${BUILD_ADDR} build --output name=${NAME}:${TAG},${BUILD_OPTION} --frontend=dockerfile.v0 --local context=${DOCKERFILE_FOLDER}   --local dockerfile=${DOCKERFILE_FOLDER}
build-arm: create
	@echo "--- build Dockerfile --"
	${BUILD} --addr ${BUILD_ADDR_ARM} build --output name=${NAME}:${TAG}-arm,${BUILD_OPTION} --frontend=dockerfile.v0 --local context=${DOCKERFILE_FOLDER}   --local dockerfile=${DOCKERFILE_FOLDER}
    
# buildctl build --output type=image,name=bookserver2:31000/karosu/nextcloud:25.0.3,push=true,registry.insecure=true --frontend=dockerfile.v0 --local context=./ --local dockerfile=./
# buildctl --addr tcp://buildkit-arm.bookserver.home:1235 \
#   build --output type=image,name=bookserver2:31000/karosu/test:0.1,push=true,registry.insecure=true \
#   --frontend=dockerfile.v0 \
#   --local context=. \
#   --local dockerfile=.