#!/bin/bash
name="bookserver2:31000/karosu/client-go"
version="0.3.0"
docker build -t $name:$version -f ./Dockerfile_arm ./
docker push $name:$version