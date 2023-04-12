---
page_title: "AWS EC2 Deployment Guide"
subcategory: ""
description: |-
  This guide explains how to deploy the BastionZero bzero agent on an AWS EC2 instance and create the policy required in order to connect
---

# AWS EC2 Deployment Guide

This guide explains how to deploy the BastionZero
[`bzero`](https://github.com/bastionzero/bzero) agent on an AWS EC2 instance and
create the policy required in order to connect. By the end of this guide, your
EC2 instance will be
[autodiscovered](https://docs.bastionzero.com/docs/deployment/installing-the-agent#autodiscovery)
as a Bzero target by BastionZero, and you will be able to connect to it using
the [`zli`](https://github.com/bastionzero/zli).

This document is split into modular sections to help guide you through the
entire Terraform configuration. 

## Before you begin

* You must be an administrator of your BastionZero organization in order to
register a target.
* Create a BastionZero [API
key](https://docs.bastionzero.com/docs/admin-guide/authorization#creating-an-api-key)
in order to configure the Terraform provider. Manage your API keys at the API
key panel found [here](https://cloud.bastionzero.com/admin/apikeys).
* Create a BastionZero [registration API
key](https://docs.bastionzero.com/docs/admin-guide/authorization#registration-api-keys)
in order to register the EC2 instance as a target in your BastionZero
organization.
* Use Terraform 1.x or higher.

This guide assumes you have basic knowledge of the AWS Terraform provider as it
is used to deploy an EC2 instance. Although this guide uses AWS to register a
target, please note that BastionZero supports other cloud providers as well. Use
this guide as a model to register your instances at other cloud providers.

## Setup

First, we setup the BastionZero Terraform provider and the [AWS Terraform
provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#authentication-and-configuration).



```terraform
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
```

-> **Note** Version constraints are included for the sake of completeness.
Please change the version constraints as you see fit. Keep in mind though this
guide is written assuming the versions defined in the Terraform snippet above.

For security purposes, we choose to configure the BastionZero provider using an
environment variable as it's more secure than hardcoding the secret in the
Terraform file itself.

Set the `BASTIONZERO_API_SECRET` environment variable to the API key's secret
that you created in the [previous step](#before-you-begin) before running
`terraform apply`.

```sh
export BASTIONZERO_API_SECRET=api-secret
```

## Create an environment

All BastionZero targets belong to a single environment. Environments help
organize a collection of targets. They're especially useful when creating
BastionZero policy when you wish to apply the same set of policy access rules to
a group of targets.

Let's create a [`bastionzero_environment`](../resources/environment) to contain
our EC2 instance once it is autodiscovered and registered as a Bzero target.

```terraform
resource "bastionzero_environment" "env" {
  name        = "example-env"
  description = "Environment managed by Terraform."
}
```

## Fetch autodiscovery script

Next, we'll query for an autodiscovery script
([`bastionzero_ad_bash`](../data-sources/ad_bash)) to use during the
provisioning of the EC2 instance.

The fetched script does not embed a registration secret which is required in
order to register the target successfully. We'll use the registration API key's
secret that you created in the ["Before you begin"](#before-you-begin) step, and
the [`replace`](https://www.terraform.io/language/functions/replace) function in
order to complete the script.

```terraform
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
```

Configure the `bzero_reg_secret` input variable by setting an [environment
variable](https://developer.hashicorp.com/terraform/language/values/variables#environment-variables)
before running `terraform apply`.

```sh
export TF_VAR_bzero_reg_secret=reg-secret
```

~> **Warning** The registration secret is sensitive data. If a malicious
attacker obtains this credential, they could register their own instances as
targets in your BastionZero organization. Once the registration secret is used
in a Terraform module (e.g. fetched via a data source), it is stored in the
Terraform state file. Please protect your state files accordingly. See
HashiCorp's article about managing sensitive data in Terraform state
[here](https://developer.hashicorp.com/terraform/language/state/sensitive-data).

## Register an EC2 instance

With a complete autodiscovery script (`local.ad_script`), we're now ready to
create the EC2 instance and have it register itself to BastionZero.

The following Terraform creates an EC2 instance running Ubuntu 20.04 in the
default VPC and creates a new security group in the default VPC in order to
reject all inbound traffic and allow all outbound traffic to/from this instance.
Please modify accordingly to use your own custom VPC and security groups for
enhanced security.

-> **Note** BastionZero supports other platforms as well. We're using Ubuntu
20.04 as an example. See
[here](https://docs.bastionzero.com/docs/deployment/installing-the-agent#step-1-agent-installation)
to view the list of supported `bzero` agent platforms.

```terraform
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
```

## Control access via policy

Your newly created EC2 instance should appear in the
["Targets"](https://cloud.bastionzero.com/admin/targets) list shortly. Now let's
create a target connect policy to control who in your organization has access to
your target.

```terraform
# Get all users and groups
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
```

See the [`bastionzero_targetconnect_policy`](../resources/targetconnect_policy)
resource documentation to learn about all the available options.

Specify the users' emails in your organization to whom you wish to give access
to in `local.users`.

If you've [synced groups](https://cloud.bastionzero.com/admin/integrations) from
your IdP (see the [SSO
Management](https://docs.bastionzero.com/docs/admin-guide/authentication/sso-management)
guide for more information), then specify the groups in your organization to
whom you wish to give access to in `local.groups`.

## Next steps

Congratulations! You've successfully registered an EC2 instance to your
BastionZero organization, and you've restricted access to it via BastionZero
policy.

If you experienced any issues with the BastionZero Terraform provider, please
open a ticket at
https://github.com/bastionzero/terraform-provider-bastionzero/issues. For
assistance on any issue with using or installing BastionZero, please reach out
to support@bastionzero.com.

Here are some suggestions on what to try next.

### Connect to your target using the `zli`

Download and install the `zli` to connect to your target. Follow the
instructions [here](https://github.com/bastionzero/zli#install).

```sh
zli login
zli connect ubuntu@instance-id
```

Replace `instance-id` with the value of the `instance_id` output.

### Setup a session recording policy

Setup a session recording policy so that shell connections to your target are
recorded. Use the
[`bastionzero_sessionrecording_policy`](../resources/sessionrecording_policy)
resource to manage the policy in Terraform. Learn more about session recording
policies
[here](https://docs.bastionzero.com/docs/admin-guide/authorization#session-recording).

### Setup a JIT policy

Enable the [Slack
integration](https://docs.bastionzero.com/docs/automation-and-integrations/slack)
in your BastionZero organization, so that you can write just-in-time (JIT)
policies to grant temporary access to your BastionZero target subject to
administrator approval.

Use the [`bastionzero_jit_policy`](../resources/jit_policy) resource to manage
the policy in Terraform. Learn more about JIT policies
[here](https://docs.bastionzero.com/docs/admin-guide/authorization#just-in-time).