data "bastionzero_dac_targets" "example" {}

# Find all DAC targets whose health endpoints are healthy
output "healthy_dacs" {
  value = [
    for each in data.bastionzero_dac_targets.example.targets
    : each if each.status == "Online"
  ]
}
