resource "aws_s3_bucket" "project_data_bucket" {
  bucket = local.bucket_name

  tags = var.common_tags

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_s3_bucket_policy" "secure_bucket_policy" {
  bucket = aws_s3_bucket.project_data_bucket.id
  policy = data.aws_iam_policy_document.secure_bucket_policy_doc.json
}

resource "aws_s3_bucket_server_side_encryption_configuration" "secure_bucket_encryption" {
  bucket = aws_s3_bucket.project_data_bucket.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.project_data_encryption.arn
    }
    bucket_key_enabled = true
  }
}

resource "aws_s3_bucket_public_access_block" "project_data_bucket" {
  bucket = aws_s3_bucket.project_data_bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_ownership_controls" "project_data_bucket_ownership" {
  bucket = aws_s3_bucket.project_data_bucket.id

  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}


data "aws_iam_policy_document" "secure_bucket_policy_doc" {
  statement {
    sid    = "EnforceSecureTransport"
    effect = "Deny"
    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
    actions   = ["s3:*"]
    resources = local.bucket_and_objects_arn
    condition {
      test     = "Bool"
      variable = "aws:SecureTransport"
      values   = ["false"]
    }
  }

  statement {
    sid    = "EnforceTLS12"
    effect = "Deny"
    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
    actions   = ["s3:*"]
    resources = local.bucket_and_objects_arn
    condition {
      test     = "NumericLessThan"
      variable = "aws:TlsVersion"
      values   = ["1.2"]
    }
  }

  statement {
    sid    = "DenyUnencryptedUploads"
    effect = "Deny"
    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
    actions   = ["s3:PutObject"]
    resources = ["${aws_s3_bucket.project_data_bucket.arn}/*"]
    condition {
      test     = "StringNotEquals"
      variable = "s3:x-amz-server-side-encryption"
      values   = ["aws:kms"]
    }
  }
}
