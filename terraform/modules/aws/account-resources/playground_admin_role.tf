data "aws_iam_policy_document" "admin_role_trust_policy" {
  statement {
    sid    = "AllowAssumeRole"
    effect = "Allow"
    principals {
      type        = "AWS"
      identifiers = [var.admin_non_root_user_arn]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "admin_role" {
  name               = "PlaygroundAdministratorRole"
  description        = "Administrative role for managing my projects"
  assume_role_policy = data.aws_iam_policy_document.admin_role_trust_policy.json

  tags = var.common_tags
}

resource "aws_iam_role_policy_attachment" "admin_access_attachment" {
  role       = aws_iam_role.admin_role.name
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}

data "aws_iam_policy_document" "assume_admin_role_policy_doc" {
  statement {
    effect    = "Allow"
    actions   = ["sts:AssumeRole"]
    resources = [aws_iam_role.admin_role.arn]
  }
}

resource "aws_iam_policy" "allow_assume_admin_role_policy" {
  name        = "AllowAssumePlaygroundAdminRolePolicy"
  description = "Allows assuming the project administrator role"
  policy      = data.aws_iam_policy_document.assume_admin_role_policy_doc.json

  tags = var.common_tags
}
