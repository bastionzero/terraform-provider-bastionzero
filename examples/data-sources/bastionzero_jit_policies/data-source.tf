data "bastionzero_jit_policies" "example" {}

# Find all JIT policies whose durations provide just in time access greater than
# 1 hour.
locals {
  large_durations = [
    for each in data.bastionzero_jit_policies.example.policies
    : each if each.duration > 60
  ]
}