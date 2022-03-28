
provider "azurerm" {
  features {}
}

module "agent" {
  source = "./modules/agent"

  for_each         = var.instance_regions
  instance_regions = each.key
  instance_count   = 10
  name             = var.name
}