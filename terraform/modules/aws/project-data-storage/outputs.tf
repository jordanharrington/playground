output "bucket" {
  value = aws_s3_bucket.project_data_bucket
}

output "bucket_key_arn" {
  value = aws_kms_key.project_data_encryption.arn
}
