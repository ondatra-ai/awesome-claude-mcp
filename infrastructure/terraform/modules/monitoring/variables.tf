variable "cluster_name" {
  description = "ECS Cluster name for dashboard metrics"
  type        = string
}

variable "log_retention_days" {
  description = "CloudWatch Logs retention period"
  type        = number
  default     = 30
}

variable "alb_arn_suffix" {
  description = "ALB ARN suffix used for CloudWatch metric dimensions (e.g., app/xyz/abc)"
  type        = string
}

variable "alarm_topic_arn" {
  description = "SNS Topic ARN for alarm notifications"
  type        = string
}
