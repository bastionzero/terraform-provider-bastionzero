# Using Terraform >= 1.x syntax

data "bastionzero_web_targets" "example" {}

# Find target with specific name. `web_target` is null if not found.
output "web_target" {
  value = one([
    for each in data.bastionzero_web_targets.example.targets
    : each if each.name == "example-target"
  ])
}
