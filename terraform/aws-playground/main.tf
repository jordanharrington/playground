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