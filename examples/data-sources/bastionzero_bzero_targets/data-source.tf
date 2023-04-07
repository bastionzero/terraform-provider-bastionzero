data "bastionzero_bzero_targets" "example" {}

# Find all Bzero targets whose names contain "ubuntu"
locals {
  ubuntu_targets = [
    for each in data.bastionzero_bzero_targets.example.targets
    : each if can(regex("ubuntu", each.name))
  ]
}
