#!/bin/bash

TEMP_CONTAINER_NAME=hasher-builder-`git rev-parse --short --verify HEAD`
BRANCH=$1
docker build -t tchaudhry/rainbow-hasher-go:scratch-builder -f Dockerfile.full .
docker run --name=$TEMP_CONTAINER_NAME tchaudhry/rainbow-hasher-go:scratch-builder 
docker cp $TEMP_CONTAINER_NAME:/app/bin ./
docker container rm $TEMP_CONTAINER_NAME
docker build -t tchaudhry/rainbow-hasher-go:$BRANCH -f Dockerfile.scratch .
rm -rf bin
docker push tchaudhry/rainbow-hasher-go:$BRANCH
