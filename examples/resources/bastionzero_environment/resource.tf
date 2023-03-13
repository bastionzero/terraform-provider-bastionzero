# Create a new environment
resource "bastionzero_environment" "example" {
  name                          = "example-env"
  offline_cleanup_timeout_hours = 12
}