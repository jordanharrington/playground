variable "common_tags" {
  type = map(string)
  description = "A map of common tags to apply to all resources."
  default = {}
}

variable "admin_non_root_user_arn" {
  type        = string
  description = "A map of common tags to apply to all resources."
}