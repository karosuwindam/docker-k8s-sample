#!/bin/bash
directory=bookmark

if [ ! -d $directory ]; then
    echo "create bookmark dir"
    mkdir $directory
fi

if [ -z "$(ls bookmark)" ]; then
    echo "copy bookmarks.html"
    cp bookmarks.html bookmark
fi

./app