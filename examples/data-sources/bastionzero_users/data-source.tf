data "bastionzero_users" "example" {}

# Create map from email address to name
output "email_to_name_map" {
  value = {
    for each in data.bastionzero_users.example.users
    : each.email => each.full_name
  }
}

# Find all admins
output "admin_users" {
  value = [
    for each in data.bastionzero_users.example.users
    : each if each.is_admin
  ]
}