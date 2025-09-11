variable "aws_region" {
  description = "AWS region to deploy resources in"
  type        = string
}

variable "environment" {
  description = "Deployment environment name (dev|staging|prod)"
  type        = string
}
