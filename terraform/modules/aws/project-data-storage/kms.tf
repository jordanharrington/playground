data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "bucket_encryption_key_policy" {
  statement {
    sid       = "AllowFullAccessForRootAndKeyAdmins"
    effect    = "Allow"
    resources = ["*"]
    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root",
        var.key_admin_arn
      ]
    }
    actions = ["kms:*"]
  }

  statement {
    sid       = "AllowUseOfTheKeyForServices"
    effect    = "Allow"
    resources = ["*"]

    principals {
      type        = "AWS"
      identifiers = [var.key_admin_arn]
    }
    actions = [
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*",
      "kms:DescribeKey",
    ]
  }

  statement {
    sid    = "AllowS3ToUseTheKeyForTheBucket"
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["s3.amazonaws.com"]
    }
    actions = [
      "kms:GenerateDataKey*",
      "kms:Encrypt",
    ]
    resources = ["*"]
    condition {
      test     = "ArnLike"
      variable = "aws:SourceArn"
      values   = [aws_s3_bucket.project_data_bucket.arn]
    }
  }
}

resource "aws_kms_key" "project_data_encryption" {
  description         = "KMS key for encrypting the project data S3 bucket"
  enable_key_rotation = true
  policy              = data.aws_iam_policy_document.bucket_encryption_key_policy.json

  tags = var.common_tags
}

resource "aws_kms_alias" "project_data_encryption_key_alias" {
  name          = "alias/${local.bucket_name}-encryption-key"
  target_key_id = aws_kms_key.project_data_encryption.id
}
