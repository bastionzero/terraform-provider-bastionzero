data "bastionzero_groups" "example" {}

# Find all groups whose names are equal to "Engineering"
locals {
  engineering_groups = [
    for each in data.bastionzero_groups.example.groups
    : each if each.name == "Engineering"
  ]
}
