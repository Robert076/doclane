variable "project" {}
variable "env" {}
variable "account_id" {}

resource "aws_s3_bucket" "this" {
  bucket = "${var.project}-${var.env}-${var.account_id}"

  tags = {
    Project     = var.project
    Environment = var.env
  }
}

resource "aws_iam_policy" "s3_policy" {
  name        = "${var.project}-${var.env}-policy"
  description = "S3 CRUD policy for ${var.env}"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:GetObjectVersion",
          "s3:DeleteObjectVersion"
        ]
        Resource = "${aws_s3_bucket.this.arn}/*"
      },
      {
        Effect = "Allow"
        Action = [
          "s3:ListBucket"
        ]
        Resource = aws_s3_bucket.this.arn
      }
    ]
  })
}

output "bucket_name" {
  value = aws_s3_bucket.this.bucket
}

output "s3_policy_arn" {
  value = aws_iam_policy.s3_policy.arn
}
