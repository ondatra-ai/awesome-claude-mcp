resource "aws_ecr_repository" "this" {
  for_each             = var.create_base_repos ? toset(var.repositories) : toset([])
  name                 = each.value
  image_tag_mutability = "MUTABLE"
  force_delete         = var.environment != "prod"
  image_scanning_configuration {
    scan_on_push = true
  }
  encryption_configuration {
    encryption_type = "AES256"
  }
}

# Create environment-specific repositories (service-env pattern)
resource "aws_ecr_repository" "env_specific" {
  for_each             = var.create_env_specific_repos ? toset([for repo in var.repositories : "${repo}-${var.environment}"]) : toset([])
  name                 = each.value
  image_tag_mutability = "MUTABLE"
  force_delete         = var.environment != "prod"
  image_scanning_configuration {
    scan_on_push = true
  }
  encryption_configuration {
    encryption_type = "AES256"
  }
}
