locals {
  gh_actions_runner_prefix = "aws-playground-github-actions-runner"
}

resource "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list = ["sts.amazonaws.com"]

  thumbprint_list = ["6938fd4d98bab03faadb97b34396831e3780aea1"]
}

data "aws_iam_policy_document" "gh_oidc_assume_role" {
  statement {
    sid     = "GitHubOIDCAssumeRole"
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.github.arn]
    }

    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = ["sts.amazonaws.com"]
    }

    condition {
      test     = "StringLike"
      variable = "token.actions.githubusercontent.com:sub"
      values = [
        for r in var.gh_runner_allowed_repos : "repo:${r}:*"
      ]
    }
  }
}

resource "aws_iam_role" "github_actions_runner_role" {
  name               = "${local.gh_actions_runner_prefix}-ecr-role"
  assume_role_policy = data.aws_iam_policy_document.gh_oidc_assume_role.json
  tags               = var.common_tags
}

data "aws_iam_policy_document" "ecr_push_policy" {
  statement {
    sid       = "ECRAuth"
    effect    = "Allow"
    actions   = ["ecr:GetAuthorizationToken"]
    resources = ["*"]
  }

  statement {
    sid    = "PushPullSpecificRepo"
    effect = "Allow"
    actions = [
      "ecr:BatchCheckLayerAvailability",
      "ecr:CompleteLayerUpload",
      "ecr:InitiateLayerUpload",
      "ecr:PutImage",
      "ecr:UploadLayerPart",
      "ecr:BatchGetImage",
      "ecr:DescribeImages",
      "ecr:DescribeRepositories",
    ]
    resources = [
      aws_ecr_repository.playground_ecr.arn
    ]
  }
}

resource "aws_iam_policy" "ecr_push_policy" {
  name   = "${local.gh_actions_runner_prefix}-ecr-push"
  policy = data.aws_iam_policy_document.ecr_push_policy.json
  tags   = var.common_tags
}

resource "aws_iam_role_policy_attachment" "attach_ecr_push" {
  role       = aws_iam_role.github_actions_runner_role.name
  policy_arn = aws_iam_policy.ecr_push_policy.arn
}