resource "aws_cloudwatch_log_group" "services" {
  for_each          = toset(["frontend", "backend", "mcp-service"])
  name              = "/aws/ecs/${each.value}"
  retention_in_days = var.log_retention_days
}
