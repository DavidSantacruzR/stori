output "ecr_repository_url" {
  description = "container registry url"
  value       = aws_ecr_repository.lambda_ecr_repo.repository_url
}
