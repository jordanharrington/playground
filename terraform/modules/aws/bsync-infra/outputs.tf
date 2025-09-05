output "ecr_repository_url" {
  value = aws_ecr_repository.bsync-gateway.repository_url
}

output "github_actions_role_arn" {
  value = aws_iam_role.github_actions_runner_role.arn
}