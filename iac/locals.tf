locals {
  tags = {
    project     = var.project
    environment = var.environment
    managed_by  = var.managed_by
    created_at  = var.created_at
  }
}
