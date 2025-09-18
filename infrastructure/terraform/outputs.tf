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
