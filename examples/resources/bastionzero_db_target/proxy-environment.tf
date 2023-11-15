data "bastionzero_environments" "example" {}

locals {
  # Find environment with name "example-env". `env` is null if not found
  env = one([
    for each in data.bastionzero_environments.example.environments
    : each if each.name == "example-env"
  ])
}

resource "bastionzero_db_target" "example" {
  name                 = "example-psql-db-target"
  remote_host          = "localhost"
  environment_id       = local.env.id
  remote_port          = 5432
  proxy_environment_id = local.env.id
  # Configures the Db daemon to run on port 5432 when a user connects to this Db
  # target
  local_port = 5432
}
