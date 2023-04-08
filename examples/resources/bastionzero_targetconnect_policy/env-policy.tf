# Get all users, groups, and environments 
data "bastionzero_users" "u" {}
data "bastionzero_groups" "g" {}
data "bastionzero_environments" "e" {}

locals {
  # Define, by email address, users to add to the policy
  users = ["alice@example.com", "bob@example.com", "charlie@example.com"]
  # Define, by name, the groups and environments to add to the policy
  groups = ["Engineering", "Product"]
  envs   = ["Default", "Demo-Env"]
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
  environments = [
    for each in data.bastionzero_environments.e.environments
    : each.id if contains(local.envs, each.name)
  ]

  # Permit access as "ec2-user"
  target_users = ["ec2-user"]
  # Allow shell access, file upload/download, and SSH
  verbs = ["Shell", "FileTransfer", "Tunnel"]
}
