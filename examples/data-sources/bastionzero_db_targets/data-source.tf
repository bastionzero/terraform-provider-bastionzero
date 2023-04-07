data "bastionzero_db_targets" "example" {}

# Find all Db targets with remote port 5432
locals {
  psql_targets = [
    for each in data.bastionzero_db_targets.example.targets
    : each if each.remote_port == 5432
  ]
}
