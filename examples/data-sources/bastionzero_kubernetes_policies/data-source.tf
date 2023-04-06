data "bastionzero_kubernetes_policies" "example" {}

# Find all Kubernetes policies whose lists of cluster groups contain
# "system:masters"
locals {
  sudo_policies = [
    for each in data.bastionzero_kubernetes_policies.example.policies
    : each if contains(each.cluster_groups, "system:masters")
  ]
}
