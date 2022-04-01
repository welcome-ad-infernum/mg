output "resource_group_list" {
  description = "Azure resource group name"
  value       = [
    for k in module.agent: k.azurerm_resource_group_name
  ]
}

output "container_list" {
  description = "Azure container name"
  value       = [
    for k in module.agent: k.azurerm_container_group_name
  ]
}