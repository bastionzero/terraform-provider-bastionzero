data "bastionzero_kubernetes_policies" "example" {}

# Find all Kubernetes policies whose lists of cluster groups contain
# "system:masters"
output "sudo_policies" {
  value = [
    for each in data.bastionzero_kubernetes_policies.example.policies
    : each if contains(each.cluster_groups, "system:masters")
  ]
}
