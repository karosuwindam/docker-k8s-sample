

docker run --rm --name buildkit -d -v /home/pi/docker:/etc/docker:z --privileged -p 1234:1234 moby/buildkit --addr tcp://0.0.0.0:1234

buildctl build --output type=image,name=bookserver2:31000/karosu/test:0.1,push=true,registry.insecure=true --frontend=dockerfile.v0 --local context=. --local dockerfile=.