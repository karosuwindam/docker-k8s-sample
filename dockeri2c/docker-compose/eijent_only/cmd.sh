#!/bin/bash
name="bookserver2:31000/moni/dockeri2cmoni"
version="0.83"
docker build -t $name:$version ./services/app
docker push $name:$version