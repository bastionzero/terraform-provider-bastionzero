data "bastionzero_environments" "example" {}

# Find all environments whose names contain "test"
output "test_envs" {
  value = [
    for each in data.bastionzero_environments.example.environments
    : each if can(regex("test", each.name))
  ]
}
