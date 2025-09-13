data "aws_iam_policy_document" "lambda_trust_policy_doc" {
  statement {
    effect = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda_exec" {
  name               = var.lambda_role_name
  assume_role_policy = data.aws_iam_policy_document.lambda_trust_policy_doc.json
}

resource "aws_iam_role_policy_attachment" "bucket_policy_attachment" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.bucket_policy.arn
}


data "aws_iam_policy_document" "logs_policy_doc" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "logs_policy" {
  name   = "${var.lambda_role_name}-logs"
  policy = data.aws_iam_policy_document.logs_policy_doc.json
}

resource "aws_iam_role_policy_attachment" "logs_policy_attachment" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.logs_policy.arn
}

resource "aws_lambda_function" "gateway" {
  function_name = "bsync-gateway"
  role          = aws_iam_role.lambda_exec.arn
  package_type  = "Image"
  image_uri     = var.lambda_image_tag
  architectures = ["x86_64"]
  memory_size   = 128
  tags          = var.common_tags
}