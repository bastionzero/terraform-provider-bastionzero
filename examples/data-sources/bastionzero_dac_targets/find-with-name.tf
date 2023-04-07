# Using Terraform >= 1.x syntax

data "bastionzero_dac_targets" "example" {}

# Find target with specific name. `dac_target` is null if not found.
locals {
  dac_target = one([
    for each in data.bastionzero_dac_targets.example.targets
    : each if each.name == "example-target"
  ])
}
