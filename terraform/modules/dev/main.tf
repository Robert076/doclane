provider "aws" {
  region = "eu-west-1"
}

module "s3" {
  source     = "../../modules/s3"
  project    = "doclane"
  env        = "dev"
  account_id = "659775407830"
}

resource "aws_iam_role" "s3_doclane_role" {
  name        = "s3-doclane-role-dev"
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
