# Get all users and groups
data "bastionzero_users" "u" {}
data "bastionzero_groups" "g" {}

locals {
  # Define, by email address, users to add to the policy
  users = ["alice@example.com", "bob@example.com", "charlie@example.com"]
  # Define, by name, the groups to add to the policy
  groups = ["Engineering"]
}

resource "bastionzero_targetconnect_policy" "example" {
  name        = "example-policy"
  description = "Policy managed by Terraform."
  subjects = [
    for each in data.bastionzero_users.u.users
    : { id = each.id, type = each.type } if contains(local.users, each.email)
  ]
  groups = [
    for each in data.bastionzero_groups.g.groups
    : { id = each.id, name = each.name } if contains(local.groups, each.name)
  ]
  # Apply this policy to the environment created earlier
  environments = [bastionzero_environment.env.id]

  # Permit access as "ubuntu"
  target_users = ["ubuntu"]
  # Allow shell access, file upload/download, and SSH
  verbs = ["Shell", "FileTransfer", "Tunnel"]
}
