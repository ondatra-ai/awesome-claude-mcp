output "redis_endpoint" {
  description = "Redis primary endpoint"
  value       = coalesce(
    aws_elasticache_replication_group.this.configuration_endpoint_address,
    aws_elasticache_replication_group.this.primary_endpoint_address
  )
}

output "security_group_id" {
  description = "Redis security group ID"
  value       = aws_security_group.redis.id
}
