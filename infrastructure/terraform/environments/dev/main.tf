terraform {
  required_version = ">= 1.6.0"
}

module "vpc" {
  source = "../../modules/vpc"
}

module "ecs" {
  source     = "../../modules/ecs"
  vpc_id     = module.vpc.vpc_id
  alb_sg_id  = module.alb.sg_id
}

module "alb" {
  source            = "../../modules/alb"
  vpc_id            = module.vpc.vpc_id
  public_subnet_ids = module.vpc.public_subnet_ids
  certificate_arn   = var.certificate_arn
}

module "ecr" {
  source = "../../modules/ecr"
}

module "redis" {
  source = "../../modules/redis"
}

module "monitoring" {
  source = "../../modules/monitoring"
}

module "iam" {
  source = "../../modules/iam"
}
