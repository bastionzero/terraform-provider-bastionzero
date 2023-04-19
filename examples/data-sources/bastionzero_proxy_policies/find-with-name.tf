# Using Terraform >= 1.x syntax

data "bastionzero_proxy_policies" "example" {}

# Find policy with specific name. `policy` is null if not found.
output "policy" {
  value = one([
    for each in data.bastionzero_proxy_policies.example.policies
    : each if each.name == "example-policy"
  ])
}
