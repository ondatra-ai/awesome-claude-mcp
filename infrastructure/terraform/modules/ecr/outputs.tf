output "repository_urls" {
  description = "Map of repository names to URLs"
  value       = { for k, r in aws_ecr_repository.this : k => r.repository_url }
}
