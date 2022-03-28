# Terraform GCP deployment

These Terraform files can be used in Docker agent deployment to GCE spot instances running Container-Optimized OS.

## To use this Terraform code you need

* [gcloud sdk](https://cloud.google.com/sdk/docs/install) to be installed on your PC
* [terraform](https://www.terraform.io/downloads) to be installed on your PC
* configured `.tfvars` file. See `example.tfvars` as an example of configuring

## Useful information

New Google Cloud Platform users receive 300 USD credits for basic usage. So, we can use these credits for agent deployment. Though GCP free plan has limitations, such as 12 vCPU total, 4 public IPs per zone and 8 vCPU per zone. 

Payed accounts can use this code on full power, providing support for more VM instances in more regions.
In example provided, we use 12 VMs in 3 different availability zones. You can customize it as you wish.

## Example of using Terraform for GCP

**Login into GCP account from CLI**

`$ gcloud auth application-default login`

**Initialize terraform modules**

`$ terraform init`

**Deploy agent**

`$ terraform apply -auto-approve -var-file example.tfvars`

## Verify if agent is running properly

1. You can list VMs in your account using

`$ gcloud compute instances list`

2. Connect to the instance via SSH

`$ ssh user@<instance-ip>` OR `$ gcloud compute ssh <instance-name>`

3. View agent's logs

`$ docker logs -f mg-agent`