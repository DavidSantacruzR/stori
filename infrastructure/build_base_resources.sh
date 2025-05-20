#!/bin/zsh

echo "Running terraform plan to check deployment configuration"
terraform plan --target=aws_ecr_repository.lambda_ecr_repo
terraform plan --target=aws_s3_bucket.lambda_bucket

echo "Running building infrastructure ..."
echo "Creating the ecr repository and s3 bucket"
terraform apply --target=aws_ecr_repository.lambda_ecr_repo --auto-approve
terraform apply --target=aws_s3_bucket.lambda_bucket --auto-approve

echo "Setting lambda iam role policy"
terraform apply --target=aws_iam_role.lambda_role --auto-approve
terraform apply --target=aws_iam_policy.lambda_s3_read --auto-approve
terraform apply --target=aws_iam_policy.allow_ses_send_raw_email --auto-approve

echo "Configure role policy"
terraform apply --target=aws_iam_role_policy_attachment.lambda_basic --auto-approve
terraform apply --target=aws_iam_role_policy_attachment.lambda_s3_read_attach --auto-approve
terraform apply --target=aws_iam_role_policy_attachment.lambda_send_raw_email --auto-approve


echo "Deployment finished"
