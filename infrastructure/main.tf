
# Initial ecr configuration
resource "aws_ecr_repository" "lambda_ecr_repo" {
  name                 = var.ecr_repository_name
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = false
  }
}

# Config s3 bucket to get the csv file from.

resource "aws_s3_bucket" "lambda_bucket" {
  bucket = "stori-challenge-david-s"
  force_destroy = true
  tags = {
    Name = "stori-challenge-david-s"
  }
}

# Role policies
resource "aws_iam_role" "lambda_role" {
  name = "stori-lambda-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_policy" "lambda_s3_read" {
  name = "stori-lambda-s3-read"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "s3:GetObject"
        ]
      }
    ]
  })
}

resource "aws_iam_policy" "allow_ses_send_raw_email" {
  name = "AllowSESSendRawEmail"
  description = "Allowing sending raw emails from lambda functions"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ses:SendRawEmail"
        ],
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_s3_read_attach" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.lambda_s3_read.arn
}

resource "aws_iam_role_policy_attachment" "lambda_send_raw_email" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.allow_ses_send_raw_email.arn
}

# Config to create the lambdas based on the images on ecr
resource "aws_lambda_function" "docker_lambda_parser" {
  function_name = "lambda-parser"
  role          = aws_iam_role.lambda_role.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.lambda_ecr_repo.repository_url}:lambda-parser"
  timeout     = var.lambda_timeout
  memory_size = var.lambda_memory_size
  environment {
    variables = var.environment_variables
  }
}

resource "aws_lambda_function" "docker_lambda_storage" {
  function_name = "lambda-storage"
  role          = aws_iam_role.lambda_role.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.lambda_ecr_repo.repository_url}:lambda-storage"
  timeout     = var.lambda_timeout
  memory_size = var.lambda_memory_size
  environment {
    variables = var.environment_variables
  }
}

resource "aws_lambda_function" "docker_lambda_summary" {
  function_name = "lambda-summary"
  role          = aws_iam_role.lambda_role.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.lambda_ecr_repo.repository_url}:lambda-summary"
  timeout     = var.lambda_timeout
  memory_size = var.lambda_memory_size
  environment {
    variables = var.environment_variables
  }
}

resource "aws_lambda_function" "docker_lambda_email" {
  function_name = "lambda-email"
  role          = aws_iam_role.lambda_role.arn
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.lambda_ecr_repo.repository_url}:lambda-email"
  timeout     = var.lambda_timeout
  memory_size = var.lambda_memory_size
  environment {
    variables = var.environment_variables
  }
}