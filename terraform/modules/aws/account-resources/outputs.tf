output "admin_role_arn" {
  value = aws_iam_role.admin_role.arn
}

output "github_actions_role_arn" {
  value = aws_iam_role.github_actions_runner_role.arn
}

output "ecr_repository_url" {
  value = aws_ecr_repository.playground_ecr.repository_url
}
