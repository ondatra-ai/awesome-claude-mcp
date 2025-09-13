variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "environment" {
  description = "Environment name (dev|staging|prod)"
  type        = string
}

variable "certificate_arn" {
  description = "ACM certificate ARN for ALB HTTPS (optional)"
  type        = string
  default     = ""
}

variable "alb_name" {
  description = "ALB name; must be unique per account/region (include env suffix)"
  type        = string
}

variable "min_count_frontend" {
  description = "Min tasks for frontend"
  type        = number
}

variable "max_count_frontend" {
  description = "Max tasks for frontend"
  type        = number
}

variable "min_count_backend" {
  description = "Min tasks for backend"
  type        = number
}

variable "max_count_backend" {
  description = "Max tasks for backend"
  type        = number
}
