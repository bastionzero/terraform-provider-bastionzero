# Get all groups and Cluster targets 
data "bastionzero_groups" "g" {}
data "bastionzero_cluster_targets" "t" {}

locals {
  # Define, by name, the groups and targets to add to the policy
  groups  = ["Administrators"]
  targets = ["prod-cluster", "stage-cluster"]
}

resource "bastionzero_kubernetes_policy" "example" {
  name        = "example-policy"
  description = "Policy managed by Terraform."
  groups = [
    for each in data.bastionzero_groups.g.groups
    : { id = each.id, name = each.name } if contains(local.groups, each.name)
  ]
  clusters = [
    for each in data.bastionzero_cluster_targets.t.targets
    : each.id if contains(local.targets, each.name)
  ]

  # Allow unrestricted rights to all Kubernetes APIs by permitting access to the
  # "system:masters" group which is a privileged group
  cluster_users  = ["cluster-admin"]
  cluster_groups = ["system:masters"]
}
