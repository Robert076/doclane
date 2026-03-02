variable "project" {}
variable "vpc_id" {}
variable "private_subnet_ids" {}
variable "db_endpoint" {}
variable "db_username" {}
variable "db_password" {}
variable "db_name" {}
variable "bucket_name" {}
variable "s3_role_arn" {}
variable "s3_policy_arn" {}
variable "jwt_secret" {}
variable "allowed_origin" {}

resource "aws_security_group" "lambda" {
  name   = "${var.project}-lambda-sg"
  vpc_id = var.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project}-lambda-sg"
  }
}

resource "aws_iam_role" "lambda_execution" {
  name = "${var.project}-lambda-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

# Allows Lambda to create network interfaces inside the VPC
resource "aws_iam_role_policy_attachment" "vpc_access" {
  role       = aws_iam_role.lambda_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# Allows Lambda to assume the S3 role
resource "aws_iam_policy" "assume_s3_role" {
  name = "${var.project}-assume-s3-role"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = "sts:AssumeRole"
        Resource = var.s3_role_arn
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "assume_s3_role_attach" {
  role       = aws_iam_role.lambda_execution.name
  policy_arn = aws_iam_policy.assume_s3_role.arn
}

resource "aws_lambda_function" "api" {
  function_name = "${var.project}-api"
  role          = aws_iam_role.lambda_execution.arn
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  handler       = "bootstrap"
  filename      = "placeholder.zip" # replaced on first real deploy

  timeout     = 30
  memory_size = 256

  vpc_config {
    subnet_ids         = var.private_subnet_ids
    security_group_ids = [aws_security_group.lambda.id]
  }

  environment {
    variables = {
      DB_HOST        = var.db_endpoint
      DB_USER        = var.db_username
      DB_PASSWORD    = var.db_password
      DB_NAME        = var.db_name
      DB_PORT        = "5432"
      JWT_SECRET     = var.jwt_secret
      S3_BUCKET_NAME = var.bucket_name
      AWS_ROLE_S3    = var.s3_role_arn
      ALLOWED_ORIGIN = var.allowed_origin
    }
  }
}
