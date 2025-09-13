locals {
  normalized_s3_buckets = [
    for t in var.target_s3_buckets : {
      bucket      = t.bucket
      region      = t.region
      kms_key_arn = try(t.kms_key_arn, null)
      ops         = t.ops
      bucket_arn  = "arn:aws:s3:::${t.bucket}"
      object_arns = (
        length(try(t.prefixes, [])) == 0
        ? ["arn:aws:s3:::${t.bucket}/*"]
        : [for p in t.prefixes : "arn:aws:s3:::${t.bucket}/${p}*"]
      )
    }
  ]

  s3_put_actions = ["s3:PutObject"]
  kms_put_actions = ["kms:Encrypt", "kms:GenerateDataKey*", "kms:DescribeKey"]

  s3_object_statements = [
    for t in local.normalized_s3_buckets : {
      actions = concat(contains(t.ops, "put") ? local.s3_put_actions : [])
      resources = t.object_arns
    }
    if length(concat(contains(t.ops, "put") ? local.s3_put_actions : [])) > 0
  ]

  kms_statements = [
    for t in local.normalized_s3_buckets : {
      region = t.region
      actions = concat(contains(t.ops, "put") ? local.kms_put_actions : [])
      resources = [t.kms_key_arn]
    }
    if try(t.kms_key_arn, null) != null
    && t.kms_key_arn != ""
    && length(concat(contains(t.ops, "put") ? local.kms_put_actions : [])) > 0
  ]
}

data "aws_iam_policy_document" "bucket_policy_doc" {
  # S3
  dynamic "statement" {
    for_each = {for idx, s in local.s3_object_statements : idx => s}
    content {
      effect    = "Allow"
      actions   = statement.value.actions
      resources = statement.value.resources
    }
  }
  # KMS
  dynamic "statement" {
    for_each = {for idx, s in local.kms_statements : idx => s}
    content {
      effect    = "Allow"
      actions   = statement.value.actions
      resources = statement.value.resources

      condition {
        test     = "StringEquals"
        variable = "kms:ViaService"
        values = ["s3.${statement.value.region}.amazonaws.com"]
      }
    }
  }
}

resource "aws_iam_policy" "bucket_policy" {
  name   = "bsync-lambda-bucket-policy"
  policy = data.aws_iam_policy_document.bucket_policy_doc.json
  tags   = var.common_tags
}