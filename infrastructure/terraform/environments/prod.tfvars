aws_region  = "us-east-1"
environment = "prod"
# alb_name constructed dynamically as "awesome-claude-mcp-alb-prod"
min_count_frontend = 2
max_count_frontend = 6
min_count_backend  = 2
max_count_backend  = 6
# certificate_arn = "arn:aws:acm:...:certificate/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
