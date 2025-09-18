output "cluster_id" {
  description = "ECS cluster ID"
  value       = aws_ecs_cluster.this.id
}

output "cluster_name" {
  description = "ECS cluster name"
  value       = aws_ecs_cluster.this.name
}

output "service_sg_id" {
  description = "ECS service security group ID"
  value       = aws_security_group.service.id
}

output "frontend_service_name" {
  description = "Frontend service name"
  value       = aws_ecs_service.frontend.name
}

output "backend_service_name" {
  description = "Backend service name"
  value       = aws_ecs_service.backend.name
}

output "backend_image_used" {
  description = "Backend image URL being used"
  value       = local.backend_image_url
}

output "frontend_image_used" {
  description = "Frontend image URL being used"
  value       = local.frontend_image_url
}
