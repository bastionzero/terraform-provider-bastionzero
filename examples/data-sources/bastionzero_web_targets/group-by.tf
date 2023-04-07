data "bastionzero_web_targets" "example" {}

# Group Web targets with same base proxy target
locals {
  web_targets_by_base = {
    for each in data.bastionzero_web_targets.example.targets
    : each.proxy_target_id => { id = each.id, name = each.name }...
  }
}
