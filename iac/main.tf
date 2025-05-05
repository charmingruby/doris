module "remote_state" {
  source = "./modules/aws/remote-state"

  tags = local.tags
}

module "database" {
  source = "./modules/aws/database"

  tags = local.tags
}
