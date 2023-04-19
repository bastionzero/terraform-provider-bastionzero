data "bastionzero_groups" "example" {}

# Find all groups whose names are equal to "Engineering"
output "engineering_groups" {
  value = [
    for each in data.bastionzero_groups.example.groups
    : each if each.name == "Engineering"
  ]
}
