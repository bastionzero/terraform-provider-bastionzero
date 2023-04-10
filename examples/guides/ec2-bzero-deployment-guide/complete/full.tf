terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
    bastionzero = {
      source  = "bastionzero/bastionzero"
      version = "~> 0.0"
    }
  }
}

# Configure the AWS provider.
provider "aws" {
  region = "us-east-1"
}

# Configure the BastionZero provider. An empty provider configuration assumes
# the BASTIONZERO_API_SECRET environment variable is set. The provider uses the
# environment variable's value as the `api_secret` field.
provider "bastionzero" {}
resource "bastionzero_environment" "env" {
  name        = "example-env"
  description = "Environment managed by Terraform."
}
# This is only an example. We recommend to fetch this secret from your preferred
# secrets manager. Do not expose a .tf file with your secret.
variable "bzero_reg_secret" {
  type        = string
  description = "BastionZero registration secret used to register a target."
  sensitive   = true
  nullable    = false
}

data "bastionzero_ad_bash" "ad_script" {
  environment_id     = bastionzero_environment.env.id
  target_name_option = "AwsEc2Metadata"
}

locals {
  ad_script = sensitive(
    replace(
      data.bastionzero_ad_bash.ad_script.script,
      "<REGISTRATION-SECRET-GOES-HERE>",
      var.bzero_reg_secret
    )
  )
}
# Query the latest ubuntu AMI for 20.04-amd64
data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

# Create security group in the default VPC
module "demo_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.5"
  name    = "demo-security-group"

  # Only permit outbound traffic. Reject all inbound traffic.
  egress_cidr_blocks = ["0.0.0.0/0"]
  egress_rules       = ["all-tcp", "all-udp", "all-icmp"]
}

# Create EC2 instance in the default VPC
module "demo_ec2_instance" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "~> 4.0"

  name          = "demo-bzero-target"
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"
  user_data     = local.ad_script

  vpc_security_group_ids = [module.demo_sg.security_group_id]

  tags = {
    Terraform = "true"
  }
}

output "instance_id" {
  value = split("instance/", module.demo_ec2_instance.arn)[1]
} # Get all users and groups
data "bastionzero_users" "u" {}
data "bastionzero_groups" "g" {}

locals {
  # Define, by email address, users to add to the policy
  users = ["alice@example.com", "bob@example.com", "charlie@example.com"]
  # Define, by name, the groups to add to the policy
  groups = ["Engineering"]
}

resource "bastionzero_targetconnect_policy" "example" {
  name        = "example-policy"
  description = "Policy managed by Terraform."
  subjects = [
    for each in data.bastionzero_users.u.users
    : { id = each.id, type = each.type } if contains(local.users, each.email)
  ]
  groups = [
    for each in data.bastionzero_groups.g.groups
    : { id = each.id, name = each.name } if contains(local.groups, each.name)
  ]
  # Apply this policy to the environment created earlier
  environments = [bastionzero_environment.env.id]

  # Permit access as "ubuntu"
  target_users = ["ubuntu"]
  # Allow shell access, file upload/download, and SSH
  verbs = ["Shell", "FileTransfer", "Tunnel"]
}
