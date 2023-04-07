data "bastionzero_bzero_targets" "example" {}

# Find all Bzero targets whose names contain "ubuntu"
output "ubuntu_targets" {
  value = [
    for each in data.bastionzero_bzero_targets.example.targets
    : each if can(regex("ubuntu", each.name))
  ]
}
