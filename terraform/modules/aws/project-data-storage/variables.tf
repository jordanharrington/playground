variable "bucket_account_id" {
  type        = string
  description = "The name of the S3 bucket."
}

variable "key_admin_arn" {
  type        = string
  description = "The ARN of the IAM principal for KMS key administration."
}

variable "common_tags" {
  type = map(string)
  description = "A map of common tags to apply to all resources."
  default = {}
}