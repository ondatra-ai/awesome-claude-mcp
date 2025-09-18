variable "cluster_name" {
  description = "ECS cluster name"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID for service security group"
  type        = string
}

variable "public_subnet_ids" {
  description = "Public subnet IDs for ECS services networking"
  type        = list(string)
}



variable "execution_role_arn" {
  description = "IAM role ARN for ECS task execution"
  type        = string
}

variable "task_role_arn" {
  description = "IAM role ARN for ECS task permissions"
  type        = string
}

variable "desired_count" {
  description = "Desired count for each service"
  type        = number
  default     = 1
}

variable "frontend_image" {
  description = "Container image for frontend service"
  type        = string
}

variable "backend_image" {
  description = "Container image for backend service"
  type        = string
}
