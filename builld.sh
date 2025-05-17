docker buildx build --build-arg LAMBDA_DIR=email --platform linux/amd64 --provenance=false -t lambda-email:v0.1.0 .
docker buildx build --build-arg LAMBDA_DIR=email --platform linux/amd64 --provenance=false -t lambda-parser:v0.1.0 .
