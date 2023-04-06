# Using Terraform >= 1.x syntax

data "bastionzero_jit_policies" "example" {}

# Find policy with specific name. `policy` is null if not found.
locals {
  policy = one([
    for each in data.bastionzero_jit_policies.example.policies
    : each if each.name == "example-policy"
  ])
}
