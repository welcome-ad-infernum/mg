module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 3.0"

  name                 = var.name
  cidr                 = "10.0.0.0/16"
  azs                  = ["${var.region}a", "${var.region}b", "${var.region}c"]
  private_subnets      = ["10.0.1.0/24"]
  public_subnets       = ["10.0.2.0/24"]
  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true

  manage_default_security_group = true
  default_security_group_name   = var.name
  default_security_group_tags = {
    "role" = "mg-agent"
  }

  default_security_group_ingress = [
    {
      cidr_blocks = "0.0.0.0/0"
      description = "Allow only ssh connect to instances"
      from_port   = 22
      protocol    = "tcp"
      self        = false
      to_port     = 22
    },
  ]

  default_security_group_egress = [
    {
      cidr_blocks = "0.0.0.0/0"
      description = "Allow all egress from instances to internete"
      from_port   = 0
      protocol    = "-1"
      self        = false
      to_port     = 0
    },
  ]

  tags = {
    "role" = "mg-agent"
  }
}