resource "azurerm_resource_group" "mg" {
  name     = "${var.instance_regions}-group"
  location = var.instance_regions
}

resource "azurerm_container_group" "mg" {
  count               = var.instance_count
  name                = "${azurerm_resource_group.mg.location}-${var.name}${count.index}"
  location            = var.instance_regions
  resource_group_name = azurerm_resource_group.mg.name
  ip_address_type     = "None"
  os_type             = "Linux"

  container {
    name   = var.name
    image  = "vladstarr/mg-agent:latest"
    cpu    = 1
    memory = 1
  }
}