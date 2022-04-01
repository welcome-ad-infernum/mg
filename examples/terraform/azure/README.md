# Terraform Azure deployment

These Terraform files can be used in Docker Agent deployment using Azure container groups.
In provided example, we deploy 10 VMs in 15 different availability zones each. You can customize it as you wish.

## To use this Terraform code you need

* [azure-cli](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) to be installed on your PC
* [terraform](https://www.terraform.io/downloads) to be installed on your PC
* configured `.tfvars` file. See `example.tfvars` as an example of configuring

## Example of using Terraform for Azure

**Login into Azure account from CLI**

`$ az login`

**Initialize terraform modules**

`$ terraform init`

**Deploy agent**

`$ terraform apply -var-file example.tfvars`

## Verify if agent is running properly

After completion of deployment, you will receive resource group names and container group names. 
You can view container's logs with:

`$ az container logs -g <resource_group_name> -n <container_group_name> --follow`