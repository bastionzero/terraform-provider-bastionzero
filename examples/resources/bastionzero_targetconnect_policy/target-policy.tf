# Get all groups and Bzero targets 
data "bastionzero_groups" "g" {}
data "bastionzero_bzero_targets" "t" {}

locals {
  # Define, by name, the groups and targets to add to the policy
  groups  = ["Product", "Marketing"]
  targets = ["demo-1", "customer-insights"]
}

resource "bastionzero_targetconnect_policy" "example" {
  name        = "example-policy"
  description = "Policy managed by Terraform."
  groups = [
    for each in data.bastionzero_groups.g.groups
    : { id = each.id, name = each.name } if contains(local.groups, each.name)
  ]
  targets = [
    for each in data.bastionzero_bzero_targets.t.targets
    : { id = each.id, type = each.type } if contains(local.targets, each.name)
  ]

  # Permit access as "ec2-user" and "demo-user"
  target_users = ["ec2-user", "demo-user"]
  # Only allow shell access
  verbs = ["Shell"]
}
