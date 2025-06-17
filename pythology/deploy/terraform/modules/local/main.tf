terraform {
  required_providers {
    minio = {
      source  = "aminueza/minio"
      version = "~> 3.5.2"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6.0"
    }
  }
}

provider "minio" {
  minio_server   = "localhost:9000"
  minio_password = var.minio_secret_key
  minio_user     = var.minio_access_key
  minio_ssl      = false
}

resource "minio_s3_bucket" "bucket" {
  bucket = "pythology-bucket"
  acl    = "public"
}