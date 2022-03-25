# Terraform AWS deployment

These Terraform files can be used to deploy AWS EC2 spot instances with pre-installed systemd agent.

## To use this Terraform code you need

* [aws-cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) to be installed on your PC
* [terraform](https://www.terraform.io/downloads) to be installed on your PC
* [generate SSH key](https://www.ssh.com/academy/ssh/keygen) or use existing, place it along with terraform files and name it as `id_rsa.pub`
* configured `.tfvars` file. See `frankfurt.tfvars` as an example of configuring
* You can also switch between different regions, using workspaces feature of Terraform

## Example of using Terraform for AWS

**Login into aws account from CLI**

`$ aws configure`

**Initialize terraform modules**

`$ terraform init`

**Create new workspace for Frankfurt region**

`$ terraform workspace new frankfurt`

**Deploy agent to Frankfurt region**

`$ terraform apply -auto-approve -var-file frankfurt.tfvars`

**Create new workspace for Singapore region**

`$ terraform workspace new singapore`

**Deploy agents to Singapore region**

`$ terraform apply -auto-approve -var-file singapore.tfvars`

## Verify if agent is running properly

After the `terraform apply` command, you will receive an outputs with VMs IP addresses. 

You can show agent's logs with 

`$ ssh ec2-user@<vm_ip> journalctl -u mg -f`