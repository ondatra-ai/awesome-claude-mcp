output "cluster_arn" {
  description = "ECS cluster ARN"
  value       = aws_ecs_cluster.this.arn
}

output "service_sg_id" {
  description = "Security group ID for ECS services"
  value       = aws_security_group.services.id
}
