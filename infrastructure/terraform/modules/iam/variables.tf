variable "execution_role_name" {
  description = "Name for ECS task execution role"
  type        = string
  default     = "ecsTaskExecutionRole"
}

variable "task_role_name" {
  description = "Name for ECS task role"
  type        = string
  default     = "ecsTaskRole"
}
