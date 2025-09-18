variable "name_prefix" {
  description = "Name prefix for ECS resources"
  type        = string
  default     = "awesome-claude-mcp"
}

variable "cluster_name" {
  description = "ECS cluster name"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID for service security group"
  type        = string
}

variable "alb_sg_id" {
  description = "ALB security group ID allowed to reach services"
  type        = string
}

variable "private_subnet_ids" {
  description = "Private subnet IDs for ECS services networking"
  type        = list(string)
}

variable "tg_frontend_arn" {
  description = "ALB target group ARN for frontend"
  type        = string
}

variable "tg_backend_arn" {
  description = "ALB target group ARN for backend"
  type        = string
}

variable "frontend_image" {
  description = "Container image for frontend service"
  type        = string
  default     = "public.ecr.aws/nginx/nginx:latest"
}

variable "backend_image" {
  description = "Container image for backend service"
  type        = string
  default     = "public.ecr.aws/amazonlinux/amazonlinux:latest"
}

variable "mcp_image" {
  description = "Container image for MCP service"
  type        = string
  default     = "public.ecr.aws/amazonlinux/amazonlinux:latest"
}

variable "execution_role_arn" {
  description = "IAM role ARN for ECS task execution"
  type        = string
}

variable "task_role_arn" {
  description = "IAM role ARN for ECS task permissions"
  type        = string
}

variable "desired_count_frontend" {
  description = "Desired count for frontend service"
  type        = number
  default     = 1
}

variable "desired_count_backend" {
  description = "Desired count for backend service"
  type        = number
  default     = 1
}

variable "desired_count_mcp" {
  description = "Desired count for MCP service"
  type        = number
  default     = 1
}

variable "min_count_frontend" {
  description = "Min tasks for frontend"
  type        = number
  default     = 1
}

variable "max_count_frontend" {
  description = "Max tasks for frontend"
  type        = number
  default     = 3
}

variable "min_count_backend" {
  description = "Min tasks for backend"
  type        = number
  default     = 1
}

variable "max_count_backend" {
  description = "Max tasks for backend"
  type        = number
  default     = 3
}

variable "scale_cpu_target" {
  description = "Target CPU utilization for scaling"
  type        = number
  default     = 60
}

variable "namespace_name" {
  description = "Private DNS namespace for service discovery"
  type        = string
  default     = "svc.local"
}
