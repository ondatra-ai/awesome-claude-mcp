resource "aws_cloudwatch_log_group" "services" {
  for_each          = toset(["frontend", "backend", "mcp-service"])
  name              = "/aws/ecs/${each.value}"
  retention_in_days = var.log_retention_days
}

resource "aws_cloudwatch_dashboard" "infra" {
  dashboard_name = "mcp-gde-infra"
  dashboard_body = jsonencode({
    widgets = [
      {
        type = "metric"
        x = 0, y = 0, width = 12, height = 6
        properties = {
          metrics = [["AWS/ApplicationELB", "HTTPCode_ELB_5XX_Count", {"stat": "Sum"}]]
          region  = "${data.aws_region.current.id}"
          title   = "ALB 5XX Errors"
        }
      },
      {
        type = "metric"
        x = 12, y = 0, width = 12, height = 6
        properties = {
          metrics = [
            ["AWS/ECS", "CPUUtilization", "ClusterName", "mcp-gde-cluster", {"stat": "Average"}],
            ["AWS/ECS", "MemoryUtilization", "ClusterName", "mcp-gde-cluster", {"stat": "Average"}]
          ]
          region = "${data.aws_region.current.id}"
          title  = "ECS Cluster Utilization"
        }
      }
    ]
  })
}

data "aws_region" "current" {}

resource "aws_cloudwatch_metric_alarm" "alb_5xx" {
  alarm_name          = "alb-5xx-errors"
  namespace           = "AWS/ApplicationELB"
  metric_name         = "HTTPCode_ELB_5XX_Count"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 1
  period              = 60
  statistic           = "Sum"
  threshold           = 1
}
