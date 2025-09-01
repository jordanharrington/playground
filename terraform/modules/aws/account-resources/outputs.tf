output "admin_role_arn" {
  value = aws_iam_role.admin_role.arn
}

output "repository_url" {
  value = aws_ecr_repository.playground_ecr.repository_url
}