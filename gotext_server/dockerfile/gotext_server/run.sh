#!/bin/bash
directory=txt

if [ ! -d $directory ]; then
    echo "create txt dir"
    mkdir $directory
fi

if [ -z "$(ls txt)" ]; then
    echo "copy txt-tmp/*"
    cp txt-tmp/* txt
fi

./app