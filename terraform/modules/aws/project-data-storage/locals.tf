locals {
  bucket_name = "project-data-storage-${var.bucket_account_id}"
  bucket_and_objects_arn = [
    aws_s3_bucket.project_data_bucket.arn,
    "${aws_s3_bucket.project_data_bucket.arn}/*",
  ]
}