aws_region  = "us-east-1"
environment = "staging"
# alb_name constructed dynamically as "awesome-claude-mcp-alb-staging"
min_count_frontend = 2
max_count_frontend = 4
min_count_backend  = 2
max_count_backend  = 4
# certificate_arn = "arn:aws:acm:...:certificate/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
