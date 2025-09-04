variable "common_tags" {
  type        = map(string)
  description = "A map of common tags to apply to all resources."
  default     = {}
}

variable "admin_non_root_user_arn" {
  type        = string
  description = "A map of common tags to apply to all resources."
}

variable "gh_runner_allowed_repos" {
  description = "List of GitHub repos allowed to assume the GitHub Actions runner role (owner/repo)."
  type        = list(string)
}