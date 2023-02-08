## 概要

Buildkitを導入して、リモートからビルドできるようにする


```
docker run --rm --name buildkit -d -v /home/pi/docker:/etc/docker:z --privileged -p 1234:1234 moby/buildkit --addr tcp://0.0.0.0:1234

export BUILDKIT_HOST=tcp://buildkit.bookserver.home:1234
export BUILDKIT_HOST=tcp://buildkit-arm.bookserver.home:1235


buildctl build --output type=image,name=bookserver2:31000/karosu/test:0.1,push=true,registry.insecure=true --frontend=dockerfile.v0 --local context=. --local dockerfile=.
buildctl build --output type=image,name=bookserver2:31000/karosu/test:0.1,push=true,registry.insecure=true --frontend=dockerfile.v0 --local context=./nginx/docker-compose/ --local dockerfile=./nginx/docker-compose/
```

## 主な変更内容について

## 参考
DockerHubの公式イメージ管理
https://github.com/moby/buildkit

以下のサイトを参考にingress controllerをカスタマイズ
https://zenn.dev/tingtt/articles/ee239a40aaca7f