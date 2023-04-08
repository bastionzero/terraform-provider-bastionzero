# Get all groups, Db targets, and Web targets
data "bastionzero_groups" "g" {}
data "bastionzero_db_targets" "t" {}
data "bastionzero_web_targets" "t" {}

locals {
  # Define, by name, the groups and targets to add to the policy
  groups  = ["Product", "Marketing"]
  targets = ["demo-psql", "grafana"]
}

resource "bastionzero_proxy_policy" "example" {
  name        = "example-policy"
  description = "Policy managed by Terraform."
  groups = [
    for each in data.bastionzero_groups.g.groups
    : { id = each.id, name = each.name } if contains(local.groups, each.name)
  ]
  targets = [
    for each in concat(
      data.bastionzero_db_targets.t.targets,
      data.bastionzero_web_targets.t.targets
    )
    : { id = each.id, type = each.type } if contains(local.targets, each.name)
  ]
}
