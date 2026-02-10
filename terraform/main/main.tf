provider "aws" {
  region = "eu-west-1"
}

module "s3" {
  source  = "../modules/s3"
  project = "doclane"
}

module "vpc" {
  source  = "../modules/vpc"
  project = "doclane"
  region  = "eu-west-1"
}

module "ecs" {
  source             = "../modules/ecs"
  project            = "doclane"
  region             = "eu-west-1"
  vpc_id             = module.vpc.vpc_id
  public_subnet_ids  = module.vpc.public_subnets_ids
  private_subnet_ids = module.vpc.private_subnets_ids
}

resource "aws_iam_role" "s3_doclane_role" {
  name        = "s3-doclane-role"
  description = "Role to be assumed when accessing the S3 Doclane bucket. Has the s3_doclane_policy attached"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          AWS = "arn:aws:iam::659775407830:user/robert-beres"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "s3_attach" {
  role       = aws_iam_role.s3_doclane_role.name
  policy_arn = module.s3.s3_policy_arn
}

output "s3_doclane_role" {
  value = aws_iam_role.s3_doclane_role.arn
}

output "vpc_id" {
  value = module.vpc.vpc_id
}

output "public_subnets_arn" {
  value = module.vpc.public_subnets_arn
}

output "private_subnets_arn" {
  value = module.vpc.private_subnets_arn
}
