terraform {
  required_version = ">= 1.6.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = var.environment
      Project     = "mcp-google-docs-editor"
    }
  }
}

module "vpc" {
  source = "../../modules/vpc"
}

module "alb" {
  source            = "../../modules/alb"
  name              = var.alb_name
  vpc_id            = module.vpc.vpc_id
  public_subnet_ids = module.vpc.public_subnet_ids
  certificate_arn   = var.certificate_arn
}

module "ecr" {
  source = "../../modules/ecr"
}

module "iam" {
  source = "../../modules/iam"
}

module "ecs" {
  source               = "../../modules/ecs"
  vpc_id               = module.vpc.vpc_id
  alb_sg_id            = module.alb.sg_id
  private_subnet_ids   = module.vpc.private_subnet_ids
  tg_frontend_arn      = module.alb.target_groups["frontend"]
  tg_backend_arn       = module.alb.target_groups["backend"]
  execution_role_arn   = module.iam.execution_role_arn
  task_role_arn        = module.iam.task_role_arn
  frontend_image       = format("%s:latest", module.ecr.repository_urls["frontend"])
  backend_image        = format("%s:latest", module.ecr.repository_urls["backend"])
  mcp_image            = format("%s:latest", module.ecr.repository_urls["mcp-service"])
  min_count_frontend   = var.min_count_frontend
  max_count_frontend   = var.max_count_frontend
  min_count_backend    = var.min_count_backend
  max_count_backend    = var.max_count_backend
}

module "redis" {
  source             = "../../modules/redis"
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  allowed_sg_ids     = [module.ecs.service_sg_id]
  name_prefix        = var.environment
}

module "monitoring" {
  source = "../../modules/monitoring"
}
