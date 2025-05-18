#!/bin/zsh

echo "Running terraform plan to check deployment configuration"
terraform plan --target=aws_ecr_repository.lambda_ecr_repo

echo "Running building infrastructure ..."
echo "Creating the ecr repository"
terraform apply --target=aws_ecr_repository.lambda_ecr_repo --auto-approve

echo "Setting lambda iam role policy"
terraform apply --target=aws_iam_role.lambda_role --auto-approve

echo "Configure role policy"
terraform apply --target=aws_iam_role_policy_attachment.lambda_basic --auto-approve


echo "Deployment finished"
