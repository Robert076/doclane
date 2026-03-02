variable "project" {}
variable "vpc_id" {}
variable "private_subnet_ids" {}
variable "lambda_sg_id" {}
variable "db_username" {}
variable "db_password" {}
variable "db_name" {}

resource "aws_db_subnet_group" "this" {
  name       = "${var.project}-db-subnet-group"
  subnet_ids = var.private_subnet_ids

  tags = {
    Name = "${var.project}-db-subnet-group"
  }
}

resource "aws_security_group" "rds" {
  name   = "${var.project}-rds-sg"
  vpc_id = var.vpc_id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [var.lambda_sg_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project}-rds-sg"
  }
}

resource "aws_db_instance" "this" {
  identifier               = "${var.project}-db"
  engine                   = "postgres"
  engine_version           = "17"
  instance_class           = "db.t4g.micro"
  storage_type             = "gp2"
  allocated_storage        = 20
  publicly_accessible      = false
  delete_automated_backups = true
  storage_encrypted        = false

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  skip_final_snapshot = true
  multi_az            = false

  tags = {
    Name = "${var.project}-db"
  }
}
