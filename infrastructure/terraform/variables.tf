variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "environment" {
  description = "Environment name (dev|staging|prod)"
  type        = string
}

variable "backend_image" {
  description = "Backend Docker image URI"
  type        = string
  default     = "nginx:alpine"
}

variable "frontend_image" {
  description = "Frontend Docker image URI"
  type        = string
  default     = "nginx:alpine"
}
