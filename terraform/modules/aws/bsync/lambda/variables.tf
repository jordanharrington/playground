variable "common_tags" {
  type = map(string)
  description = "A map of common tags to apply to all resources."
  default = {}
}

variable "lambda_role_name" {
  description = "IAM role name for the Lambda execution role"
  type        = string
}

variable "lambda_image_tag" {
  description = "Image Tag for the Lambda function"
  type        = string
}

variable "target_s3_buckets" {
  description = <<EOT
Per-bucket permissions:
- bucket: S3 bucket name
- region: AWS region for the bucket (used for KMS ViaService)
- kms_key_arn: optional CMK ARN if bucket uses SSE-KMS
- prefixes: list of allowed key prefixes ("" or empty list => all keys)
- ops: set of allowed operations: put
EOT
  type = list(object({
    bucket = string
    region = string
    kms_key_arn = optional(string)
    prefixes = optional(list(string), [])
    ops = set(string)
  }))
}
