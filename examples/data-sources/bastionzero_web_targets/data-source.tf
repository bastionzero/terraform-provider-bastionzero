data "bastionzero_web_targets" "example" {}

# Find all Web targets with remote port 80
locals {
  port_80_targets = [
    for each in data.bastionzero_web_targets.example.targets
    : each if each.remote_port == 80
  ]
}
