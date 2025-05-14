module "remote_state" {
  source = "./modules/aws/remote-state"
  tags   = local.tags
}

module "database" {
  source = "./modules/aws/database"
  tags   = local.tags
}

module "email" {
  source           = "./modules/aws/email"
  sender_email     = var.sender_email
  recipient_emails = var.recipient_emails
  tags             = local.tags
}

module "storage" {
  source      = "./modules/aws/storage"
  bucket_name = var.bucket_name
  account_arn = data.aws_caller_identity.current.arn
  tags        = local.tags
}
