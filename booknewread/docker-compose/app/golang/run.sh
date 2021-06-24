#!/bin/bash
if [ -z "$(ls $directory)" ]; then
    echo "copy bookmarks.html"
    cp bookmarks.html bookmark
fi

./app