terraform {
  required_version = ">= 1.6.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }

  backend "s3" {}
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = var.environment
      Project     = "awesome-claude-mcp"
    }
  }
}

locals {
  name_prefix  = "awesome-claude-mcp"
  cluster_name = "${local.name_prefix}-cluster-${var.environment}"
}

module "vpc" {
  source      = "./modules/vpc"
  name_prefix = local.name_prefix
}


module "ecr" {
  source                    = "./modules/ecr"
  environment               = var.environment
  create_base_repos         = false
  create_env_specific_repos = true
}

module "iam" {
  source = "./modules/iam"
}

module "ecs_simple" {
  source             = "./modules/ecs-simple"
  cluster_name       = local.cluster_name
  vpc_id             = module.vpc.vpc_id
  public_subnet_ids  = module.vpc.public_subnet_ids
  execution_role_arn = module.iam.execution_role_arn
  task_role_arn      = module.iam.task_role_arn
  frontend_image     = var.frontend_image
  backend_image      = var.backend_image
}
