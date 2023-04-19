resource "bastionzero_environment" "example" {
  name                          = "example-env"
  description                   = "Example environment"
  offline_cleanup_timeout_hours = 12
}
