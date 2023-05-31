#!/bin/bash
name="bookserver2:31000/karosu/booknewread"
version="0.16.2s"
mkdir app/golang/bookmark
docker build -t $name:$version -f ./app/Dockerfile_arm ./app
docker push $name:$version