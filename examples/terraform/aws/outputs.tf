output "ec2_instance_id" {
  description = "The ID of the instance"
  value       = "${module.ec2_spot_instance.*.id}"
}

output "ec2_instance_public_ip" {
  description = "Instance SSH connection string"
  value       = "${module.ec2_spot_instance.*.public_ip}"
}