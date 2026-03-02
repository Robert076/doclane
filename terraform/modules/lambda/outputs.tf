output "lambda_sg_id" {
  value = aws_security_group.lambda.id
}

output "lambda_function_arn" {
  value = aws_lambda_function.api.arn
}

output "lambda_function_name" {
  value = aws_lambda_function.api.function_name
}

output "lambda_execution_role_arn" {
  value = aws_iam_role.lambda_execution.arn
}
