variable "project" {}

resource "aws_s3_bucket" "this" {
  bucket = var.project

  tags = {
    Project = var.project
  }
}
