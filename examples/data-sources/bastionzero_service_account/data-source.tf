data "bastionzero_service_account" "example" {
  id = "<service-account-id>"
}

# Output this service account's JWKS URL 
output "example_env_targets" {
  value = data.bastionzero_service_account.example.jwks_url
}