output "log_group_names" {
  description = "Created CloudWatch log group names"
  value       = [for lg in aws_cloudwatch_log_group.services : lg.name]
}
