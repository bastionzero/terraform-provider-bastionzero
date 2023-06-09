data "bastionzero_web_targets" "example" {}

# Group Web targets with same proxy target
output "web_targets_by_proxy_target" {
  value = {
    for each in data.bastionzero_web_targets.example.targets
    : each.proxy_target_id => { id = each.id, name = each.name }...
  }
}
