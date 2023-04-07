data "bastionzero_dac_targets" "example" {}

# Find all DAC targets whose health endpoints are healthy
locals {
  healthy_dacs = [
    for each in data.bastionzero_dac_targets.example.targets
    : each if each.status == "Online"
  ]
}
