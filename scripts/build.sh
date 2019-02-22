#!/usr/bin/env bash

set -e

source $(dirname $0)/base.sh

# Build the container image.
# This will generate an image with both "latest" tag and another one with
# the latest git commit.
docker build -t $NAME .
docker tag $NAME $IMG:latest
