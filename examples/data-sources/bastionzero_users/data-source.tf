data "bastionzero_users" "example" {}

locals {
  # Create map from email address to name
  email_to_name_map = {
    for each in data.bastionzero_users.example.users
    : each.email => each.full_name
  }

  # Find all admins
  admin_users = [
    for each in data.bastionzero_users.example.users
    : each if each.is_admin
  ]
}