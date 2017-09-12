#!/bin/bash

TEMP_CONTAINER_NAME=hasher-builder-`git rev-parse --short --verify HEAD`
docker build -t sha256-rainbow-hasher:scratch-builder -f Dockerfile.full .
docker run --name=$TEMP_CONTAINER_NAME sha256-rainbow-hasher:scratch-builder 
docker cp $TEMP_CONTAINER_NAME:/app/bin ./
docker container rm $TEMP_CONTAINER_NAME
docker build -t sha256-rainbow-hasher:master -f Dockerfile.scratch .
rm -rf bin
