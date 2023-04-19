# Get all Bzero targets 
data "bastionzero_bzero_targets" "t" {}

locals {
  # Define, by name, the targets to add to the policy
  bzero_targets = ["prod-target"]
}

resource "bastionzero_targetconnect_policy" "child1" {
  name        = "child-tc-policy"
  description = "Child policy managed by Terraform."
  targets = [
    for each in data.bastionzero_bzero_targets.t.targets
    : { id = each.id, type = each.type } if contains(local.bzero_targets, each.name)
  ]

  # Notice, no subjects or groups are specified. We want to limit who has root
  # access via a just-in-time policy.

  # Permit access as the "root" user
  target_users = ["root"]
  verbs        = ["Shell", "FileTransfer", "Tunnel"]
}

# Get all Cluster targets
data "bastionzero_cluster_targets" "t" {}

locals {
  # Define, by name, the targets to add to the policy
  cluster_targets = ["prod-cluster"]
}

resource "bastionzero_kubernetes_policy" "child2" {
  name        = "child-kube-policy"
  description = "Child policy managed by Terraform."
  clusters = [
    for each in data.bastionzero_cluster_targets.t.targets
    : each.id if contains(local.cluster_targets, each.name)
  ]

  # Notice, no subjects or groups are specified. We want to limit who has
  # privileged access via a just-in-time policy.

  # Allow unrestricted rights to all Kubernetes APIs by permitting access to the
  # "system:masters" group which is a privileged group
  cluster_users  = ["cluster-admin"]
  cluster_groups = ["system:masters"]
}

# Get all users and groups
data "bastionzero_users" "u" {}
data "bastionzero_groups" "g" {}

locals {
  # Define, by email address, users to add to the policy
  users = ["alice@example.com", "bob@example.com", "charlie@example.com"]
  # Define, by name, the groups to add to the policy
  groups = ["Engineering"]
}

resource "bastionzero_jit_policy" "example" {
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

  # Reference the policies created above
  child_policies = [
    for each in [
      bastionzero_targetconnect_policy.child1,
      bastionzero_kubernetes_policy.child2
    ]
    : { id = each.id }
  ]

  # Require approval by an admin
  auto_approved = false
  # Allow access only for 30 minutes
  duration = 30
}
