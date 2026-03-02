provider "aws" {
  region = var.region
}

module "vpc" {
  source         = "../modules/vpc"
  project        = var.project_name
  region         = var.region
  private_subnet = var.private_subnet_list
  public_subnet  = var.public_subnet_list
}

module "rds" {
  source             = "../modules/rds"
  db_username        = var.db_username
  db_password        = var.db_password
  db_name            = var.project_name
  project            = var.project_name
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = var.private_subnet_list
}

module "s3" {
  source  = "../modules/s3"
  project = var.project_name
}

module "frontend" {
  source  = "../modules/frontend"
  project = var.project_name
}

module "lambda" {
  source             = "../modules/lambda"
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = var.private_subnet_list
  s3_policy_arn      = module.s3.s3_policy_arn
  s3_role_arn        = module.s3.s3_doclane_role_arn
  db_username        = var.db_username
  db_password        = var.db_password
  db_name            = var.project_name
  project            = var.project_name
  bucket_name        = module.s3.bucket_name
  db_endpoint        = module.rds.rds_endpoint
  jwt_secret         = var.jwt_secret
  allowed_origin     = "https://${module.frontend.cloudfront_domain_name}"
}

resource "aws_security_group_rule" "allow_lambda_to_rds" {
  type      = "ingress"
  from_port = 5432
  to_port   = 5432
  protocol  = "tcp"

  security_group_id = module.rds.rds_sg_id

  source_security_group_id = module.lambda.lambda_sg_id
}
