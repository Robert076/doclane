terraform {
  required_version = ">= 1.5"

  backend "s3" {
    bucket         = "doclane-tfstate-659775407830"
    key            = "workload/terraform.tfstate"
    region         = "eu-west-1"
    dynamodb_table = "doclane-tfstate-locks"
    encrypt        = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.35"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.17"
    }
  }
}

provider "aws" {
  region = "eu-west-1"
}


data "terraform_remote_state" "compute" {
  backend = "s3"
  config = {
    bucket = "doclane-tfstate-659775407830"
    key    = "platform-compute/terraform.tfstate"
    region = "eu-west-1"
  }
}


data "terraform_remote_state" "data" {
  backend = "s3"
  config = {
    bucket = "doclane-tfstate-659775407830"
    key    = "platform-data/terraform.tfstate"
    region = "eu-west-1"
  }
}


data "aws_eks_cluster_auth" "main" {
  name = data.terraform_remote_state.compute.outputs.eks_cluster_name
}

provider "kubernetes" {
  host                   = data.terraform_remote_state.compute.outputs.eks_cluster_endpoint
  cluster_ca_certificate = base64decode(data.terraform_remote_state.compute.outputs.eks_cluster_ca)
  token                  = data.aws_eks_cluster_auth.main.token
}

provider "helm" {
  kubernetes {
    host                   = data.terraform_remote_state.compute.outputs.eks_cluster_endpoint
    cluster_ca_certificate = base64decode(data.terraform_remote_state.compute.outputs.eks_cluster_ca)
    token                  = data.aws_eks_cluster_auth.main.token
  }
}

locals {
  cluster_name            = data.terraform_remote_state.compute.outputs.eks_cluster_name
  backend_pod_role_arn    = data.terraform_remote_state.compute.outputs.backend_pod_role_arn
  alb_controller_role_arn = data.terraform_remote_state.compute.outputs.alb_controller_role_arn
  rds_address             = data.terraform_remote_state.compute.outputs.rds_address
  s3_bucket               = data.terraform_remote_state.compute.outputs.s3_documents_bucket
  cognito_pool_id         = data.terraform_remote_state.data.outputs.cognito_user_pool_id
  cognito_client_id       = data.terraform_remote_state.data.outputs.cognito_prod_client_id
  ecr_backend             = data.terraform_remote_state.data.outputs.ecr_backend_url
  ecr_frontend            = data.terraform_remote_state.data.outputs.ecr_frontend_url
}
