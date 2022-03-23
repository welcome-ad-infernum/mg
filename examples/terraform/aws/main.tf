provider "aws" {
  region = var.region
}

data "aws_ami" "al2" {
  owners      = ["amazon"]
  most_recent = true
  filter {
    name   = "name"
    values = ["amzn2-ami-kernel-*-hvm-*-x86_64-gp2"]
  }
}

module "key_pair_external" {
  source = "terraform-aws-modules/key-pair/aws"

  key_name   = var.name
  public_key = file("${path.root}/id_rsa.pub")

  tags = {
    role = "mg-agent"
  }
}

module "ec2_spot_instance" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "~> 3.0"

  create_spot_instance = true
  putin_khuylo         = true

  count = var.instance_count
  name  = "${var.name}-${count.index}"
  ami   = data.aws_ami.al2.id

  spot_instance_interruption_behavior = "terminate"
  spot_wait_for_fulfillment           = true
  instance_type                       = var.instance_type

  vpc_security_group_ids      = [module.vpc.default_security_group_id]
  subnet_id                   = element(module.vpc.public_subnets, 0)
  associate_public_ip_address = true

  key_name = module.key_pair_external.key_pair_key_name

  user_data_base64 = base64encode(<<EOF
#!/bin/bash -xe
exec > >(tee /var/log/user-data.log|logger -t user-data -s 2>/dev/console) 2>&1
    yum -y update
    amazon-linux-extras install docker
    systemctl start docker
    systemctl enable docker
    usermod -aG docker ec2-user
    docker run -d --restart always vladstarr/mg-agent:latest
EOF
  )

  tags = {
    role = "mg-agent"
  }
}