data "bastionzero_supported_database_configs" "example" {}

# Dictionary of all supported database authentication configurations. Keyed by
# the config's label.
output "supported_configs_map" {
  value = { for c in data.bastionzero_supported_database_configs.example.configs : c.label => c }
}
