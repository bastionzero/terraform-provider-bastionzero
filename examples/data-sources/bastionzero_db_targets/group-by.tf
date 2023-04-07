data "bastionzero_db_targets" "example" {}

# Group Db targets with same base proxy target
locals {
  db_targets_by_base = {
    for each in data.bastionzero_db_targets.example.targets
    : each.proxy_target_id => { id = each.id, name = each.name }...
  }
}
