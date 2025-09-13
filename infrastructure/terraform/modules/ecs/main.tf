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

  # Allow intra-services traffic to MCP on 9090
  ingress {
    description = "Services -> MCP"
    from_port   = 9090
    to_port     = 9090
    protocol    = "tcp"
    self        = true
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Service discovery namespace (private DNS)
resource "aws_service_discovery_private_dns_namespace" "this" {
  name        = var.namespace_name
  description = "Private namespace for ECS services"
  vpc         = var.vpc_id
}

resource "aws_service_discovery_service" "backend" {
  name = "backend"
  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.this.id
    dns_records {
      ttl  = 10
      type = "A"
    }
    routing_policy = "MULTIVALUE"
  }
  health_check_custom_config {}
}

locals {
  log_group_frontend = "/aws/ecs/frontend"
  log_group_backend  = "/aws/ecs/backend"
  log_group_mcp      = "/aws/ecs/mcp-service"
}

resource "aws_ecs_task_definition" "frontend" {
  family                   = "frontend"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = var.execution_role_arn
  task_role_arn            = var.task_role_arn
  container_definitions = jsonencode([
    {
      name      = "frontend"
      image     = var.frontend_image
      essential = true
      portMappings = [{
        containerPort = 3000
        hostPort      = 3000
        protocol      = "tcp"
      }]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = local.log_group_frontend
          awslogs-region        = data.aws_region.current.id
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])
}

resource "aws_ecs_task_definition" "backend" {
  family                   = "backend"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = var.execution_role_arn
  task_role_arn            = var.task_role_arn
  container_definitions = jsonencode([
    {
      name      = "backend"
      image     = var.backend_image
      essential = true
      portMappings = [{
        containerPort = 8080
        hostPort      = 8080
        protocol      = "tcp"
      }]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = local.log_group_backend
          awslogs-region        = data.aws_region.current.id
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])
}

resource "aws_ecs_task_definition" "mcp" {
  family                   = "mcp-service"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = var.execution_role_arn
  task_role_arn            = var.task_role_arn
  container_definitions = jsonencode([
    {
      name      = "mcp"
      image     = var.mcp_image
      essential = true
      portMappings = [{
        containerPort = 9090
        hostPort      = 9090
        protocol      = "tcp"
      }]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = local.log_group_mcp
          awslogs-region        = data.aws_region.current.id
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])
}

data "aws_region" "current" {}

resource "aws_ecs_service" "frontend" {
  name            = "frontend"
  cluster         = aws_ecs_cluster.this.arn
  task_definition = aws_ecs_task_definition.frontend.arn
  desired_count   = var.desired_count_frontend
  launch_type     = "FARGATE"
  platform_version = "LATEST"
  network_configuration {
    subnets         = var.private_subnet_ids
    security_groups = [aws_security_group.services.id]
    assign_public_ip = false
  }
  load_balancer {
    target_group_arn = var.tg_frontend_arn
    container_name   = "frontend"
    container_port   = 3000
  }
}

resource "aws_ecs_service" "backend" {
  name            = "backend"
  cluster         = aws_ecs_cluster.this.arn
  task_definition = aws_ecs_task_definition.backend.arn
  desired_count   = var.desired_count_backend
  launch_type     = "FARGATE"
  platform_version = "LATEST"
  network_configuration {
    subnets         = var.private_subnet_ids
    security_groups = [aws_security_group.services.id]
    assign_public_ip = false
  }
  load_balancer {
    target_group_arn = var.tg_backend_arn
    container_name   = "backend"
    container_port   = 8080
  }
  service_registries {
    registry_arn = aws_service_discovery_service.backend.arn
  }
}

resource "aws_ecs_service" "mcp" {
  name            = "mcp-service"
  cluster         = aws_ecs_cluster.this.arn
  task_definition = aws_ecs_task_definition.mcp.arn
  desired_count   = var.desired_count_mcp
  launch_type     = "FARGATE"
  platform_version = "LATEST"
  network_configuration {
    subnets         = var.private_subnet_ids
    security_groups = [aws_security_group.services.id]
    assign_public_ip = false
  }
}

# Autoscaling targets and policies (frontend)
resource "aws_appautoscaling_target" "frontend" {
  max_capacity       = var.max_count_frontend
  min_capacity       = var.min_count_frontend
  resource_id        = "service/${aws_ecs_cluster.this.name}/${aws_ecs_service.frontend.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "frontend_cpu" {
  name               = "frontend-cpu-tt"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.frontend.resource_id
  scalable_dimension = aws_appautoscaling_target.frontend.scalable_dimension
  service_namespace  = aws_appautoscaling_target.frontend.service_namespace
  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }
    target_value       = var.scale_cpu_target
    scale_in_cooldown  = 60
    scale_out_cooldown = 60
  }
}

# Autoscaling targets and policies (backend)
resource "aws_appautoscaling_target" "backend" {
  max_capacity       = var.max_count_backend
  min_capacity       = var.min_count_backend
  resource_id        = "service/${aws_ecs_cluster.this.name}/${aws_ecs_service.backend.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "backend_cpu" {
  name               = "backend-cpu-tt"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.backend.resource_id
  scalable_dimension = aws_appautoscaling_target.backend.scalable_dimension
  service_namespace  = aws_appautoscaling_target.backend.service_namespace
  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }
    target_value       = var.scale_cpu_target
    scale_in_cooldown  = 60
    scale_out_cooldown = 60
  }
}
