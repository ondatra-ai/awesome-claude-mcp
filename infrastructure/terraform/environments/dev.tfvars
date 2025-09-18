aws_region      = "us-east-1"
environment     = "dev"
certificate_arn = "arn:aws:acm:us-east-1:195062990486:certificate/ead16455-3f73-4b99-9224-88563bd7ac17"

# Domain configuration
domain_name    = "dev.ondatra-ai.xyz"
hosted_zone_id = "Z0074068BDBQGTSF473D"

# ALB configuration (name constructed dynamically as awesome-claude-mcp-alb-dev)

# ECS scaling configuration
min_count_frontend = 1
max_count_frontend = 3
min_count_backend  = 1
max_count_backend  = 3
