output "ecr_backend_url" {
  value = aws_ecr_repository.backend.repository_url
}

output "ecr_frontend_url" {
  value = aws_ecr_repository.frontend.repository_url
}

output "gha_ecr_push_role_arn" {
  value = aws_iam_role.gha_ecr_push.arn
}

output "cognito_user_pool_id" {
  value = aws_cognito_user_pool.main.id
}

output "cognito_prod_client_id" {
  value = aws_cognito_user_pool_client.prod.id
}

output "cognito_dev_client_id" {
  value = aws_cognito_user_pool_client.dev.id
}

output "acm_cert_arn" {
  value = aws_acm_certificate.main.arn
}

output "route53_zone_id" {
  value = data.aws_route53_zone.main.zone_id
}
