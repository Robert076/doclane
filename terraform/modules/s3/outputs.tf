output "bucket_name" {
  value = aws_s3_bucket.this.bucket
}

output "s3_policy_arn" {
  value = aws_iam_policy.s3_policy.arn
}

output "s3_doclane_role" {
  value = aws_iam_role.s3_doclane_role
}

output "s3_doclane_role_arn" {
  value = aws_iam_role.s3_doclane_role.arn
}
