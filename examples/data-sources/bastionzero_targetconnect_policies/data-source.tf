data "bastionzero_targetconnect_policies" "example" {}

# Find all target connect policies whose lists of target users contain "root"
locals {
  sudo_policies = [
    for each in data.bastionzero_targetconnect_policies.example.policies
    : each if contains(each.target_users, "root")
  ]
}
