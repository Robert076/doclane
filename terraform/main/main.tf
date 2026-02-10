provider "aws" {
  region = var.region
}

module "s3" {
  source  = "../modules/s3"
  project = var.project_name
}

module "vpc" {
  source         = "../modules/vpc"
  project        = var.project_name
  region         = var.region
  private_subnet = var.private_subnet_list
  public_subnet  = var.public_subnet_list
}
