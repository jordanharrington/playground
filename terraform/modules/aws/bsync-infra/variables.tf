variable "common_tags" {
  type        = map(string)
  description = "A map of common tags to apply to all resources."
  default     = {}
}

variable "admin_non_root_user_arn" {
  type        = string
  description = "A map of common tags to apply to all resources."
}

variable "ecr_kms_key_arn" {
  type        = string
  description = "A KMS key used to encrypt Playground repositories"
}