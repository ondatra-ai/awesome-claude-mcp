output "alb_arn" {
  description = "Application Load Balancer ARN"
  value       = aws_lb.this.arn
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
