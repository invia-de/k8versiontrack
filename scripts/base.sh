#!/usr/bin/env bash
IMG=invia/k8versiontrack
TAG=$(git log -1 --pretty=%H)
NAME=$IMG:$TAG
