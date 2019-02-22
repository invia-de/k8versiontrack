#!/usr/bin/env bash

# This script is meant to be us
set -e

source $(dirname $0)/base.sh

# What this script does:
# 1. Authenticate with the registry
# 2. Push the container image with all tags (hashed and latest)
# 3. Deauths from the registry
# If no registry credentials are given, the script will exit with error code 1.

if [[ -z "$REGISTRY_USER" || -z "$REGISTRY_PASS" ]];
then
  echo "Can't authenticate with the registry. REGISTRY_USER, REGISTRY_PASS or REGISTRY_URL missing."
  exit 1
fi

docker login -u $REGISTRY_USER -p $REGISTRY_PASS

docker push $IMG

docker logout
