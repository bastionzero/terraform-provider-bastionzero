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
}