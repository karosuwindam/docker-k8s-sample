#!/bin/bash
name="bookserver2:31000/karosu/booknewread"
version="0.7.2"
docker build -t $name:$version -f ./app/Dockerfile_arm ./app
docker push $name:$version