output "ecr_repository_urls" {
  description = "Map of ECR repository names to URLs"
  value       = module.ecr.repository_urls
}

output "alb_dns_name" {
  description = "ALB DNS name for the dev environment"
  value       = module.alb.alb_dns_name
}
