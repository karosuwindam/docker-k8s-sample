#!/bin/bash
name="bookserver2:31000/moni/slackout"
version="0.4"
docker build -t $name:$version .
docker push $name:$version
