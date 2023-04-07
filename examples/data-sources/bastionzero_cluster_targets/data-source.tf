data "bastionzero_cluster_targets" "example" {}

# Create set of valid cluster users across all cluster targets
output "valid_cluster_users" {
  value = toset(flatten([
    for each in data.bastionzero_cluster_targets.example.targets
    : each.valid_cluster_users
  ]))
}
