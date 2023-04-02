bookserver2:31000/karosu/nextcloud:25.0.3

export BUILDKIT_HOST=tcp://buildkit.bookserver.home:1234

buildctl build --output type=image,name=bookserver2:31000/karosu/nextcloud:25.0.3,push=true,registry.insecure=true --frontend=dockerfile.v0 --local context=./ --local dockerfile=./