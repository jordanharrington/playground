module "account-resources" {
  source                  = "../modules/aws/account-resources"
  admin_non_root_user_arn = var.playground_admin_user_arn
  common_tags             = local.common_tags
}

module "project-data-storage" {
  source            = "../modules/aws/project-data-storage"
  bucket_account_id = var.playground_account_id
  key_admin_arn     = module.account-resources.admin_role_arn
  common_tags       = local.common_tags
}

module "bsync-infra" {
  source                  = "../modules/aws/bsync/infra"
  admin_non_root_user_arn = var.playground_admin_user_arn
  ecr_kms_key_arn         = module.account-resources.ecr_kms_key_arn
  common_tags             = local.common_tags
}

module "bsync-lambda" {
  source           = "../modules/aws/bsync/lambda"
  lambda_image_tag = "${module.bsync-infra.ecr_repository_url}:1.0.0"
  lambda_role_name = "bsync-lambda-exec"
  target_s3_buckets = [
    {
      bucket      = module.project-data-storage.bucket.bucket
      region      = module.project-data-storage.bucket.region
      kms_key_arn = module.project-data-storage.bucket_key_arn
      ops = toset(["put"])
    },
  ]
  common_tags = local.common_tags
}

module "bsync-aws-gateway" {
  source      = "../modules/aws/bsync/gateway"
  name        = "bsync-gateway"
  lambda_arn  = module.bsync-lambda.lambda_arn
  common_tags = local.common_tags
}