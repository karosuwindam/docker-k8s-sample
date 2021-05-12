#!/bin/sh

name="bookserver2:31000/moni/slackpost"
version="0.1"

if test $2;then
    op=$2
fi

case $1 in
    "run")
        docker run --rm -d -p 4500:8080 --name=slackpost $name:$version
    ;;
    "down")
        docker stop slackpost
    ;;
    "build")
        docker build -t $name:$version ./
    ;;
    "push")
        docker push $name:$version
    ;;
    "help"|*)
        echo "./cmd.sh [build|push|help]"
        exit
    ;;
esac