# Get all users and service accounts
data "bastionzero_users" "u" {}
data "bastionzero_service_accounts" "u" {}

resource "bastionzero_sessionrecording_policy" "example" {
  name        = "example-policy"
  description = "Policy managed by Terraform."
  subjects = [
    for each in concat(
      data.bastionzero_users.u.users,
      data.bastionzero_service_accounts.u.service_accounts
    )
    : { id = each.id, type = each.type }
  ]
}
