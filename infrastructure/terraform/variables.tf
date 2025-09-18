variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "environment" {
  description = "Environment name (dev|staging|prod)"
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
  validation {
    condition     = var.min_count_backend >= 0
    error_message = "min_count_backend must be >= 0."
  }
}

variable "max_count_backend" {
  description = "Max tasks for backend"
  type        = number
  validation {
    condition     = var.max_count_backend >= 1
    error_message = "max_count_backend must be >= 1."
  }
}

variable "certificate_arn" {
  description = "ARN of the SSL certificate for HTTPS (optional - ALB will only support HTTP if not provided)"
  type        = string
  default     = null
}

variable "domain_name" {
  description = "Custom domain name for the application (e.g., dev.ondatra-ai.xyz)"
  type        = string
  default     = null
}

variable "hosted_zone_id" {
  description = "Route53 hosted zone ID for the domain (required if domain_name is provided)"
  type        = string
  default     = null
}
