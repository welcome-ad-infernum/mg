output "azurerm_resource_group_name" {
  description = "Azure resource group name"
  value       = azurerm_resource_group.mg.name
}

output "azurerm_container_group_name" {
  description = "Azure container group name"
  value       = "${azurerm_container_group.mg.*.name}"
}