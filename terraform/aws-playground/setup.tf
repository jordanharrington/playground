locals {
  common_tags = {
    source_project    = "https://github.com/jordanharrington/playground"
    terraform_managed = "true"
  }
}

terraform {
  required_version = "~> 1.13.1"
  backend "s3" {
    bucket       = "playground-tf-state-81920374"
    key          = "playground/terraform.tfstate"
    use_lockfile = true
    region       = "us-east-1"
    encrypt      = true
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.10.0"
    }
  }
}