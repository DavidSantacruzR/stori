#!/bin/zsh

echo "Deploying lambdas to cloud provider"
terraform apply --target=aws_lambda_function.docker_lambda_parser --auto-approve
terraform apply --target=aws_lambda_function.docker_lambda_storage --auto-approve
terraform apply --target=aws_lambda_function.docker_lambda_summary --auto-approve
terraform apply --target=aws_lambda_function.docker_lambda_email --auto-approve

echo "Configuring lambda step functions orchestration"
terraform apply --target=aws_iam_role.step_function_role --auto-approve
terraform apply --target=aws_iam_role_policy.step_function_policy --auto-approve
terraform apply --target=aws_sfn_state_machine.lambda_pipeline --auto-approve