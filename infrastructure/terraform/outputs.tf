output "environment" {
  description = "Current environment"
  value       = var.environment
}

output "region" {
  description = "AWS region used"
  value       = var.aws_region
}

output "cluster_name" {
  description = "ECS cluster name"
  value       = module.ecs_simple.cluster_name
}

output "ecr_repositories" {
  description = "ECR repository URLs"
  value       = module.ecr.repository_urls
}

output "ecr_registry_id" {
  description = "ECR registry ID (AWS Account ID)"
  value       = module.ecr.registry_id
}

output "backend_repository_url" {
  description = "Backend ECR repository URL"
  value       = module.ecr.backend_repository_url
}

output "frontend_repository_url" {
  description = "Frontend ECR repository URL"
  value       = module.ecr.frontend_repository_url
}

output "backend_image_used" {
  description = "Backend image URL actually used in ECS task definition"
  value       = module.ecs_simple.backend_image_used
}

output "frontend_image_used" {
  description = "Frontend image URL actually used in ECS task definition"
  value       = module.ecs_simple.frontend_image_used
}
