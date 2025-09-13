# Configure a remote backend for team collaboration (uncomment and customize)
# terraform {
#   backend "s3" {
#     bucket         = "your-tfstate-bucket"
#     key            = "mcp-google-docs-editor/${var.environment}/terraform.tfstate"
#     region         = var.aws_region
#     dynamodb_table = "your-tfstate-locks"
#     encrypt        = true
#   }
# }
