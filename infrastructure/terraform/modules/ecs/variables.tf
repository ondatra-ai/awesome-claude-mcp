variable "cluster_name" {
  description = "ECS cluster name"
  type        = string
  default     = "mcp-gde-cluster"
}

variable "vpc_id" {
  description = "VPC ID for service security group"
  type        = string
}

variable "alb_sg_id" {
  description = "ALB security group ID allowed to reach services"
  type        = string
}
