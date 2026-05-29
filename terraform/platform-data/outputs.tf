output "ecr_backend_url" {
  value = aws_ecr_repository.backend.repository_url
}

output "ecr_frontend_url" {
  value = aws_ecr_repository.frontend.repository_url
}

output "gha_ecr_push_role_arn" {
  value = aws_iam_role.gha_ecr_push.arn
}
