# Using Terraform >= 1.x syntax

data "bastionzero_db_targets" "example" {}

# Find target with specific name. `db_target` is null if not found.
output "db_target" {
  value = one([
    for each in data.bastionzero_db_targets.example.targets
    : each if each.name == "example-target"
  ])
}
