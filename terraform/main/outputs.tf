output "documents_s3_bucket_name" {
  value = module.s3.bucket_name
}

output "s3_doclane_role" {
  value = module.s3.s3_doclane_role
}

output "vpc_id" {
  value = module.vpc.vpc_id
}

output "public_subnet_arn" {
  value = module.vpc.public_subnet_arn
}

output "public_subnet_id" {
  value = module.vpc.public_subnet_id
}

output "private_subnet_arn" {
  value = module.vpc.private_subnet_arn
}

output "private_subnet_id" {
  value = module.vpc.public_subnet_id
}
