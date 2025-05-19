#!/bin/zsh

echo "Deploying lambdas to cloud provider"
terraform apply --target=aws_lambda_function.docker_lambda_parser --auto-approve
terraform apply --target=aws_lambda_function.docker_lambda_storage --auto-approve
terraform apply --target=aws_lambda_function.docker_lambda_summary --auto-approve
terraform apply --target=aws_lambda_function.docker_lambda_email --auto-approve
