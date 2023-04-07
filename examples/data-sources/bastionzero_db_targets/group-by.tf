data "bastionzero_db_targets" "example" {}

# Group Db targets with same base proxy target
output "db_targets_by_base" {
  value = {
    for each in data.bastionzero_db_targets.example.targets
    : each.proxy_target_id => { id = each.id, name = each.name }...
  }
}
