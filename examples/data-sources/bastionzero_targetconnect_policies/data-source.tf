data "bastionzero_targetconnect_policies" "example" {}

# Find all target connect policies whose lists of target users contain "root"
output "sudo_policies" {
  value = [
    for each in data.bastionzero_targetconnect_policies.example.policies
    : each if contains(each.target_users, "root")
  ]
}
