# Using Terraform >= 1.x syntax

data "bastionzero_web_targets" "example" {}

# Find target with specific name. `web_target` is null if not found.
locals {
  web_target = one([
    for each in data.bastionzero_web_targets.example.targets
    : each if each.name == "example-target"
  ])
}
