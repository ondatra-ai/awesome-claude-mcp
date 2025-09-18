output "repository_urls" {
  description = "Map of repository names to URLs"
  value = merge(
    { for k, r in aws_ecr_repository.this : k => r.repository_url },
    { for k, r in aws_ecr_repository.env_specific : k => r.repository_url }
  )
}

output "registry_id" {
  description = "ECR registry ID"
  value       = data.aws_caller_identity.current.account_id
}

output "backend_repository_url" {
  description = "Backend repository URL"
  value       = var.create_env_specific_repos ? aws_ecr_repository.env_specific["backend-${var.environment}"].repository_url : (var.create_base_repos ? aws_ecr_repository.this["backend"].repository_url : null)
}

output "frontend_repository_url" {
  description = "Frontend repository URL"
  value       = var.create_env_specific_repos ? aws_ecr_repository.env_specific["frontend-${var.environment}"].repository_url : (var.create_base_repos ? aws_ecr_repository.this["frontend"].repository_url : null)
}
