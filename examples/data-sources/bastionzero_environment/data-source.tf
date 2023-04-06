data "bastionzero_environment" "example" {
  id = "<environment-id>"
}

# Output this environment's targets
output "example_env_targets" {
  value = data.bastionzero_environment.example.targets
}