data "bastionzero_web_targets" "example" {}

# Find all Web targets with remote port 80
output "port_80_targets" {
  value = [
    for each in data.bastionzero_web_targets.example.targets
    : each if each.remote_port == 80
  ]
}
