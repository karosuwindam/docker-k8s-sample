#!/bin/bash
name="bookserver2:31000/moni/co2moni"
version="0.2"
docker build -t $name:$version ./
docker push $name:$version