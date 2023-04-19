data "bastionzero_service_accounts" "example" {}

# Find all service accounts whose JWKS URLs contain "google"
output "google_sas" {
  value = [
    for each in data.bastionzero_service_accounts.example.service_accounts
    : each if can(regex("google", each.jwks_url))
  ]
}
