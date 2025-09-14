variable "name" {
  type = string
}

variable "lambda_arn" {
  type = string
}

variable "common_tags" {
  type = map(string)
  description = "A map of common tags to apply to all resources."
  default = {}
}