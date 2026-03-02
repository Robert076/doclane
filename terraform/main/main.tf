provider "aws" {
  region = var.region
}

module "rds" {
  source = "../modules/rds"
  db_username = var.db_username
  db_password = var.db_password
  db_name = var.project_name
  project = var.project_name
  vpc_id = module.vpc.vpc_id
  private_subnet_ids = var.private_subnet_list
  lambda_sg_id = module
}

module "s3" {
  source  = "../modules/s3"
  project = var.project_name
  lambda_execution_role_arn = 
}

module "vpc" {
  source         = "../modules/vpc"
  project        = var.project_name
  region         = var.region
  private_subnet = var.private_subnet_list
  public_subnet  = var.public_subnet_list
}

module "lambda" {
  source = "../modules/lambda"
  vpc_id = module.vpc.vpc_id
  private_subnet_ids = var.private_subnet_list
  s3_policy_arn = module.s3.s3_policy_arn
  s3_role_arn = module.s3.s3_doclane_role_arn
  db_username = var.db_username
  db_password = var.db_password
  db_name = var.project_name
  project = var.project_name
  bucket_name = module.s3.bucket_name
}