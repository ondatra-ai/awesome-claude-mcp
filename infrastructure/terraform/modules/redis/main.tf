resource "aws_security_group" "redis" {
  name        = "${var.name_prefix}-redis-sg"
  description = "Allow access to Redis from ECS services"
  vpc_id      = var.vpc_id

  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = var.allowed_sg_ids
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_elasticache_subnet_group" "this" {
  name       = "${var.name_prefix}-redis-subnets"
  subnet_ids = var.private_subnet_ids
}

resource "aws_elasticache_replication_group" "this" {
  replication_group_id          = "${var.name_prefix}-redis-rg"
  description                   = "Redis replication group"
  engine                        = "redis"
  engine_version                = var.engine_version
  node_type                     = var.node_type
  port                          = 6379
  automatic_failover_enabled    = true
  multi_az_enabled              = true
  cluster_mode {
    num_node_groups         = var.cluster_num_node_groups
    replicas_per_node_group = var.cluster_replicas_per_node_group
  }
  at_rest_encryption_enabled    = true
  transit_encryption_enabled    = true
  subnet_group_name             = aws_elasticache_subnet_group.this.name
  security_group_ids            = [aws_security_group.redis.id]
}
