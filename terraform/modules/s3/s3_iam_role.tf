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
  policy_arn = aws_iam_policy.s3_policy.arn
}
