data "bastionzero_service_accounts" "example" {}

# Find all service accounts whose JWKS URLs contain "google"
locals {
  google_sas = [
    for each in data.bastionzero_service_accounts.example.service_accounts
    : each if can(regex("google", each.jwks_url))
  ]
}
