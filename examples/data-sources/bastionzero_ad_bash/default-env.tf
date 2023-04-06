# Using Terraform >= 1.x syntax

# Get all environments
data "bastionzero_environments" "envs" {}

# Find the default environment which is guaranteed to exist
locals {
  default_env_id = one([
    for each in data.bastionzero_environments.envs.environments
    : each if each.is_default
  ]).id
}

data "bastionzero_ad_bash" "example" {
  environment_id     = local.default_env_id
  target_name_option = "BashHostName"
}
