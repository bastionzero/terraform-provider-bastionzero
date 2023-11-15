data "bastionzero_environments" "example" {}
data "bastionzero_bzero_targets" "example" {}

locals {
  # Find environment with name "example-env". `env` is null if not found
  env = one([
    for each in data.bastionzero_environments.example.environments
    : each if each.name == "example-env"
  ])
  # Find Linux or Windows target with name "ubuntu". `proxy_target` is null if
  # not found
  proxy_target = one([
    for each in data.bastionzero_bzero_targets.example.targets
    : each if each.name == "ubuntu"
  ])
}

resource "bastionzero_db_target" "example" {
  name = "example-gcp-db-target"
  # GCP protocol prefix is required
  remote_host    = "gcp://se-demo-pwdb:us-west2:gcp-postgres"
  environment_id = local.env.id
  # Remote port has no effect in this example but it is still required
  remote_port     = 0
  proxy_target_id = local.proxy_target.id
  database_authentication_config = {
    authentication_type    = "ServiceAccountInjection"
    cloud_service_provider = "GCP"
    database               = "Postgres"
    label                  = "GCP Postgres"
  }
}
