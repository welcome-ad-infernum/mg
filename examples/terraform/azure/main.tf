
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "mg" {
  name     = "${var.region}-group"
  location = var.region
}

resource "azurerm_container_group" "mg" {
  count               = var.instance_count
  name                = "${var.region}-${var.name}${count.index}"
  location            = var.region
  resource_group_name = azurerm_resource_group.mg.name
  ip_address_type     = "None"
  os_type             = "Linux"

  container {
    name   = var.name
    image  = "vladstarr/mg-agent:latest"
    cpu    = 1
    memory = 1
  }

  tags = {
    role = var.name
  }
}