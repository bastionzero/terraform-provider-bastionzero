# Create cluster rolebinding to the built-in `view` ClusterRole
resource "kubernetes_cluster_role_binding" "viewer_cluster_role_binding" {
  metadata {
    name = "viewer-cluster-rolebinding"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    # A default cluster role that comes built-in with Kubernetes installations.
    # It allows read-access to non-sensitive information. See
    # https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles
    # for a full description.
    name = "view"
  }
  subject {
    kind = "User"
    # The cluster username that should be specified in BastionZero Kubernetes
    # policy (`cluster_users`)
    name      = "viewer"
    api_group = "rbac.authorization.k8s.io"
  }
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

resource "bastionzero_kubernetes_policy" "example" {
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

  # Permit access as the "viewer" Kubernetes RBAC user, and therefore assume all
  # the permissions granted to the "viewer" user as defined by in-cluster
  # RoleBindings and ClusterRoleBindings.
  #
  # We created a ClusterRoleBinding for this user using the Kubernetes provider
  # above.
  cluster_users = ["viewer"]
}
