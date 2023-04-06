# Using Terraform >= 1.x syntax

data "bastionzero_environments" "example" {}

# Find environment with specific name. `environment` is null if not found.
locals {
  environment = one([
    for each in data.bastionzero_environments.example.environments
    : each if each.name == "example-env"
  ])
}
