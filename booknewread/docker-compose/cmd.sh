#!/bin/bash
name="bookserver2:31000/karosu/booknewread"
version="0.2"
docker build -t $name:$version -f ./golang/Dockerfile_arm ./golang
docker push $name:$version