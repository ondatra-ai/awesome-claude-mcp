output "alb_arn" {
  description = "Application Load Balancer ARN"
  value       = aws_lb.this.arn
}

output "alb_arn_suffix" {
  description = "ALB ARN suffix for CloudWatch metrics"
  value       = aws_lb.this.arn_suffix
}

output "alb_dns_name" {
  description = "ALB DNS name"
  value       = aws_lb.this.dns_name
}

output "target_groups" {
  description = "Target group ARNs"
  value = {
    frontend = aws_lb_target_group.frontend.arn
    backend  = aws_lb_target_group.backend.arn
  }
}

output "sg_id" {
  description = "ALB security group ID"
  value       = aws_security_group.alb.id
}

output "domain_name" {
  description = "Custom domain name (if configured)"
  value       = var.domain_name
}

output "route53_record_fqdn" {
  description = "Route53 record FQDN (if domain is configured)"
  value       = var.domain_name != null ? aws_route53_record.alb[0].fqdn : null
}
