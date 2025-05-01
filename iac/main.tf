module "remote_state" {
  source = "./modules/aws/remote-state"

  tags = local.tags
}
