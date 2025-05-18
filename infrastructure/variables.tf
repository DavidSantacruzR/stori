variable "aws_region" {
  default = "us-east-1"
}

variable "ecr_repository_name" {
  description = "ecr description"
  type        = string
  default     = "stori"
}

variable "lambda_timeout" {
  description = "lambda timeout"
  type        = number
  default     = 30
}

variable "lambda_memory_size" {
  description = "lambda mem size"
  type        = number
  default     = 256
}

variable "environment_variables" {
  description = "lambda env variables"
  type        = map(string)
  default     = {
    ENVIRONMENT = "dev"
  }
}