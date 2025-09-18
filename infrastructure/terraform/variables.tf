variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "environment" {
  description = "Environment name (dev|staging|prod)"
  type        = string
}

variable "backend_image" {
  description = "Backend Docker image URI (required)"
  type        = string
}

variable "frontend_image" {
  description = "Frontend Docker image URI (required)"
  type        = string
}
