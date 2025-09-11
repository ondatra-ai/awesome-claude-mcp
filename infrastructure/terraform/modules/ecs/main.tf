resource "aws_ecs_cluster" "this" {
  name = var.cluster_name
  setting {
    name  = "containerInsights"
    value = "enabled"
  }
  configuration {
    execute_command_configuration {
      logging = "DEFAULT"
    }
  }
  tags = {
    Name = var.cluster_name
  }
}

resource "aws_ecs_cluster_capacity_providers" "this" {
  cluster_name       = aws_ecs_cluster.this.name
  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  default_capacity_provider_strategy {
    base              = 1
    weight            = 1
    capacity_provider = "FARGATE"
  }
}

resource "aws_security_group" "services" {
  name        = "${var.cluster_name}-services-sg"
  description = "Allow ALB to reach services"
  vpc_id      = var.vpc_id

  ingress {
    description     = "Frontend from ALB"
    from_port       = 3000
    to_port         = 3000
    protocol        = "tcp"
    security_groups = [var.alb_sg_id]
  }

  ingress {
    description     = "Backend from ALB"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [var.alb_sg_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
