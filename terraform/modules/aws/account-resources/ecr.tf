data "aws_caller_identity" "current" {}

resource "random_id" "repo_suffix" {
  byte_length = 4
}

data "aws_iam_policy_document" "ecr_kms_key_policy_doc" {
  statement {
    sid       = "AllowFullAccessForRootAndKeyAdmins"
    effect    = "Allow"
    actions   = ["kms:*"]
    resources = ["*"]
    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root",
        var.admin_non_root_user_arn
      ]
    }
  }

  statement {
    sid    = "AllowServiceAccess"
    effect = "Allow"
    actions = [
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:DescribeKey"
    ]
    resources = ["*"]
    principals {
      type        = "Service"
      identifiers = ["ecr.amazonaws.com"]
    }
  }
}

resource "aws_kms_key" "playground_ecr_kms_key" {
  description         = "KMS key for encrypting the Playground ECR"
  enable_key_rotation = true
  policy              = data.aws_iam_policy_document.ecr_kms_key_policy_doc.json

  tags = var.common_tags
}

resource "aws_kms_alias" "playground_ecr_key_alias" {
  name          = "alias/${aws_ecr_repository.playground_ecr.name}-encryption-key"
  target_key_id = aws_kms_key.playground_ecr_kms_key.id
}

data "aws_iam_policy_document" "playground_ecr_policy_doc" {
  statement {
    sid    = "AllowAdminAndRootAccess"
    effect = "Allow"
    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root",
        var.admin_non_root_user_arn
      ]
    }
    actions = ["ecr:*"]
  }
}

resource "aws_ecr_repository" "playground_ecr" {
  name = "aws-playground-${random_id.repo_suffix.hex}"

  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "KMS"
    kms_key         = aws_kms_key.playground_ecr_kms_key.arn
  }

  tags = var.common_tags
}

resource "aws_ecr_repository_policy" "playground_ecr_policy" {
  repository = aws_ecr_repository.playground_ecr.name
  policy     = data.aws_iam_policy_document.playground_ecr_policy_doc.json
}

data "aws_ecr_lifecycle_policy_document" "playground_ecr_lifecycle_policy_document" {
  rule {
    priority    = 1
    description = "Expire untagged images older than 7 days"
    selection {
      count_unit   = "days"
      count_number = 7
      count_type   = "sinceImagePushed"
      tag_status   = "untagged"
    }
    action {
      type = "expire"
    }
  }
}

resource "aws_ecr_lifecycle_policy" "playground_ecr_lifecycle_policy" {
  repository = aws_ecr_repository.playground_ecr.name
  policy     = data.aws_ecr_lifecycle_policy_document.playground_ecr_lifecycle_policy_document.json
}