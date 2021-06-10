#!/bin/bash
name="bookserver2:31000/karosu/client-go"
version="0.0.1"
docker build -t $name:$version -f ./app/Dockerfile_arm ./app
docker push $name:$version