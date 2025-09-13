terraform {
  required_version = ">= 1.6.0"
}

module "vpc" {
  source = "../../modules/vpc"
}

module "ecs" {
  source     = "../../modules/ecs"
  vpc_id            = module.vpc.vpc_id
  alb_sg_id         = module.alb.sg_id
  private_subnet_ids = module.vpc.private_subnet_ids
  tg_frontend_arn    = module.alb.target_groups["frontend"]
  tg_backend_arn     = module.alb.target_groups["backend"]
  execution_role_arn = module.iam.execution_role_arn
  task_role_arn      = module.iam.task_role_arn
  min_count_frontend = 2
  max_count_frontend = 4
  min_count_backend  = 2
  max_count_backend  = 4
}

module "alb" {
  source            = "../../modules/alb"
  name              = "mcp-gde-alb-staging"
  vpc_id            = module.vpc.vpc_id
  public_subnet_ids = module.vpc.public_subnet_ids
  certificate_arn   = var.certificate_arn
}

module "ecr" {
  source = "../../modules/ecr"
}

module "redis" {
  source             = "../../modules/redis"
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  allowed_sg_ids     = [module.ecs.service_sg_id]
}

module "monitoring" {
  source = "../../modules/monitoring"
}
