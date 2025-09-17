aws_region     = "us-east-1"
environment    = "dev"
certificate_arn = "arn:aws:acm:us-east-1:195062990486:certificate/ead16455-3f73-4b99-9224-88563bd7ac17"

# ALB configuration
alb_name = "mcp-docs-dev-alb"

# ECS scaling configuration
min_count_frontend = 1
max_count_frontend = 3
min_count_backend  = 1
max_count_backend  = 3
