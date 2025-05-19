#!/bin/zsh
export $(cat .env | xargs)

PROJECT_NAME="stori"

aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $REGISTRY_URL

# lambda deploy for the parser
docker tag lambda-parser $REGISTRY_URL/$PROJECT_NAME:lambda-parser
docker push $REGISTRY_URL/$PROJECT_NAME:lambda-parser

# lambda deploy for the storage
docker tag lambda-storage $REGISTRY_URL/$PROJECT_NAME:lambda-storage
docker push $REGISTRY_URL/$PROJECT_NAME:lambda-storage

# lambda deploy for the summary
docker tag lambda-summary $REGISTRY_URL/$PROJECT_NAME:lambda-summary
docker push $REGISTRY_URL/$PROJECT_NAME:lambda-summary

# lambda deploy for the email stuff
docker tag lambda-summary $REGISTRY_URL/$PROJECT_NAME:lambda-email
docker push $REGISTRY_URL/$PROJECT_NAME:lambda-email
