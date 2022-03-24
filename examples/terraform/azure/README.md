# Terraform Azure deployment

These Terraform files can be used to deploy Microsoft Azure container groups with Docker agent.

## To use this Terraform code you need

* [azure-cli](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) to be installed on your PC
* [terraform](https://www.terraform.io/downloads) to be installed on your PC
* configured `.tfvars` file. See `example.tfvars` as an example of configuring
* You can also switch between different regions, using workspaces feature of Terraform

## Example of using Terraform for Azure

**Login into Azure account from CLI**

`$ az login`

**Initialize terraform modules**

`$ terraform init`

**Create new workspace for East Asia region**

`$ terraform workspace new eastasia`

**Deploy agent to East Asia region**

`$ terraform apply -auto-approve -var-file eastasia.tfvars`

**Create new workspace for North Europe region**

`$ terraform workspace new northeurope`

**Deploy agents to North Europe region**

`$ terraform apply -auto-approve -var-file northeurope.tfvars`

## Verify if agent is running properly

After completion of deployment, you will receive resource group name and container group names. 
You can view container's logs with:

`$ az container logs -g <resource_group_name> -n <container_group_name> --follow`