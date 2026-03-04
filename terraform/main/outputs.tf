output "documents_s3_bucket_name" {
  value = module.s3.bucket_name
}

output "s3_doclane_role" {
  value = module.s3.s3_doclane_role
}
