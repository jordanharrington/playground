data "aws_caller_identity" "current" {}

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
  name          = "alias/playground-ecr-encryption-key"
  target_key_id = aws_kms_key.playground_ecr_kms_key.id
}