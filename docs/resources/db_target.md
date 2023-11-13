---
page_title: "bastionzero_db_target Resource - terraform-provider-bastionzero"
subcategory: "Target"
description: |-
  Provides a BastionZero database target. Database targets configure remote access to database servers running on Bzero bzero_target targets or Cluster cluster_target targets.
---

# bastionzero_db_target (Resource)

Provides a BastionZero database target. Database targets configure remote access to database servers running on [Bzero](bzero_target) targets or [Cluster](cluster_target) targets.

Learn more about Db targets [here](https://docs.bastionzero.com/docs/deployment/installing-the-agent#databases).

~> **Note on proxy target/environment** A Db target _must_ be configured with
either a [`proxy_target_id`](#proxy_target_id) or a
[`proxy_environment_id`](#proxy_environment_id).

~> **Warning** Be aware that _editing_ a Db target's
[`proxy_target_id`](#proxy_target_id) or
[`proxy_environment_id`](#proxy_environment_id) after initial resource creation
may interrupt active user workflows; please consider closing all connections to
the target before modifying either of these attributes.

### Database authentication configuration

A Db target's
[`database_authentication_config`](#database_authentication_config) attribute is
optional. If left unconfigured, then the default, non-passwordless database
configuration is provided on your behalf.

-> **Note** If you don't know what combination of values to use for
[`database_authentication_config`](#database_authentication_config), then you
can use the
[`bastionzero_supported_database_configs`](supported_database_configs) data
source to get a list of supported values.

~> **Warning** Be aware that _editing_ a Db target's
[`database_authentication_config.authentication_type`](#authentication_type)
after initial resource creation may require that the Db target is reconfigured
as needed for the new authentication type. For example, if you change the
attribute to `ServiceAccountInjection` and
[`database_authentication_config.cloud_service_provider`](#cloud_service_provider)
to `GCP`, then you must also update the target's [`remote_host`](#remote_host)
to include the `gcp://` protocol prefix.

~> **Warning** Be aware that _editing_ a Db target's
[`database_authentication_config.database`](#database) after initial resource
creation may require you to reconfigure the target's allowed target users as
governed by BastionZero policy.

## Example Usage

### Db target via proxy target

Create a Db target with the default authentication configuration
(non-passwordless), and use a Bzero target to proxy the connection to the
configured database.

```terraform
data "bastionzero_environments" "example" {}
data "bastionzero_bzero_targets" "example" {}

locals {
  # Find environment with name "example-env". `env` is null if not found
  env = one([
    for each in data.bastionzero_environments.example.environments
    : each if each.name == "example-env"
  ])
  # Find Bzero target with name "ubuntu". `proxy_target` is null if not found
  proxy_target = one([
    for each in data.bastionzero_bzero_targets.example.targets
    : each if each.name == "ubuntu"
  ])
}

resource "bastionzero_db_target" "example" {
  name            = "example-psql-db-target"
  remote_host     = "localhost"
  environment_id  = local.env.id
  remote_port     = 5432
  proxy_target_id = local.proxy_target.id
}
```

### Db target via proxy environment

Create a Db target with the default authentication configuration
(non-passwordless) and a proxy environment. When a user connects to this target,
the Bzero or Cluster target with the least number of open connections in this
environment is used to proxy the connection to the configured database.

```terraform
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
```

### SplitCert Db target

```terraform
data "bastionzero_environments" "example" {}
data "bastionzero_bzero_targets" "example" {}

locals {
  # Find environment with name "example-env". `env` is null if not found
  env = one([
    for each in data.bastionzero_environments.example.environments
    : each if each.name == "example-env"
  ])
  # Find Bzero target with name "ubuntu". `proxy_target` is null if not found
  proxy_target = one([
    for each in data.bastionzero_bzero_targets.example.targets
    : each if each.name == "ubuntu"
  ])
}

resource "bastionzero_db_target" "example" {
  name            = "example-splitcert-db-target"
  remote_host     = "localhost"
  environment_id  = local.env.id
  remote_port     = 5432
  proxy_target_id = local.proxy_target.id
  database_authentication_config = {
    authentication_type = "SplitCert"
    database            = "Postgres"
    label               = "Unmanaged Postgres"
  }
}
```

### Passwordless access to Postgres on GCP Cloud

```terraform
data "bastionzero_environments" "example" {}
data "bastionzero_bzero_targets" "example" {}

locals {
  # Find environment with name "example-env". `env` is null if not found
  env = one([
    for each in data.bastionzero_environments.example.environments
    : each if each.name == "example-env"
  ])
  # Find Bzero target with name "ubuntu". `proxy_target` is null if not found
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The target's environment's ID.
- `name` (String) The target's name.
- `remote_host` (String) The target's hostname or IP address.
- `remote_port` (Number) The port of the Db server accessible via the target. This field is required for all databases; however, if `database_authentication_config.cloud_service_provider` is equal to `GCP`, then the value will be ignored when connecting to the database.

### Optional

- `database_authentication_config` (Attributes) Information about the db target's database authentication configuration. If this attribute is left unconfigured, the target is configured with the default, non-passwordless database configuration. (see [below for nested schema](#nestedatt--database_authentication_config))
- `local_port` (Number) The port of the Db daemon's localhost server that is spawned on the user's machine on connect. If this attribute is left unconfigured, an available port will be chosen when the target is connected to.
- `proxy_environment_id` (String) The target's proxy environment's ID (ID of the backing proxy environment).
- `proxy_target_id` (String) The target's proxy target's ID (ID of a [Bzero](bzero_target) or [Cluster](cluster_target) target).

### Read-Only

- `agent_public_key` (String) The target's proxy agent's public key.
- `agent_version` (String) The target's proxy agent's version.
- `id` (String) The target's unique ID.
- `last_agent_update` (String) The time this target's proxy agent last had a transition change in status formatted as a UTC timestamp string in [RFC 3339](https://datatracker.ietf.org/doc/html/rfc3339) format. Null if there has not been a single transition change.
- `region` (String) The BastionZero region that this target has connected to (follows same naming convention as AWS regions).
- `status` (String) The target's status (one of `NotActivated`, `Offline`, `Online`, `Terminated`, `Error`, or `Restarting`).
- `type` (String) The target's type (constant value `Db`).

<a id="nestedatt--database_authentication_config"></a>
### Nested Schema for `database_authentication_config`

Required:

- `authentication_type` (String) The type of authentication used when connecting to the database (one of `Default`, `SplitCert`, or `ServiceAccountInjection`).

Optional:

- `cloud_service_provider` (String) Cloud service provider hosting the database (one of `AWS`, or `GCP`). Only used for certain types of authentication (`authentication_type`), such as `ServiceAccountInjection`.
- `database` (String) The type of database running on the target (one of `CockroachDB`, `MicrosoftSQLServer`, `MongoDB`, `MySQL`, or `Postgres`).
- `label` (String) User-friendly label for this database authentication configuration.

## Import

Import is supported using the following syntax:

```shell
# A Db target can be imported by specifying the unique identifier.
terraform import bastionzero_db_target.example "01d7a020-2bbc-4ac0-b886-ac9e445d8ab1"
```