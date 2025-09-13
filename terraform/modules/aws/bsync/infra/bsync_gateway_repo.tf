data "aws_caller_identity" "current" {}

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

resource "aws_ecr_repository" "bsync-gateway" {
  name = "bsync-gateway"

  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "KMS"
    kms_key         = var.ecr_kms_key_arn
  }

  tags = var.common_tags
}

resource "aws_ecr_repository_policy" "playground_ecr_policy" {
  repository = aws_ecr_repository.bsync-gateway.name
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
  repository = aws_ecr_repository.bsync-gateway.name
  policy     = data.aws_ecr_lifecycle_policy_document.playground_ecr_lifecycle_policy_document.json
}