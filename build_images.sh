#!/bin/zsh
docker buildx build --no-cache --build-arg LAMBDA_DIR=parser --platform linux/amd64 --provenance=false -t lambda-parser .
docker buildx build --no-cache --build-arg LAMBDA_DIR=storage --platform linux/amd64 --provenance=false -t lambda-storage .
docker buildx build --no-cache --build-arg LAMBDA_DIR=summary --platform linux/amd64 --provenance=false -t lambda-summary .
docker buildx build --no-cache --build-arg LAMBDA_DIR=email --platform linux/amd64 --provenance=false -t lambda-email .