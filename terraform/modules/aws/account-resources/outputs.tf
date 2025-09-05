output "admin_role_arn" {
  value = aws_iam_role.admin_role.arn
}

output "ecr_kms_key_arn" {
  value = aws_kms_key.playground_ecr_kms_key.arn
}
