---
page_title: "bastionzero_environments Data Source - terraform-provider-bastionzero"
subcategory: "Environment"
description: |-
  Get a list of all environments in your BastionZero organization. An environment is a collection of targets.
---

# bastionzero_environments (Data Source)

Get a list of all environments in your BastionZero organization. An environment is a collection of targets.

This data source is useful if the environments in question are not managed by
Terraform, or you need to utilize any of the environments' data.

-> **Note** You can use the [`bastionzero_environment`](environment) data source
to obtain metadata about a single environment if you already know the `id`.

## Example Usage

### Basic example

```terraform
data "bastionzero_environments" "example" {}

# Find all environments whose names contain "test"
output "test_envs" {
  value = [
    for each in data.bastionzero_environments.example.environments
    : each if can(regex("test", each.name))
  ]
}
```

### Get the environment by name

```terraform
# Using Terraform >= 1.x syntax

data "bastionzero_environments" "example" {}

# Find environment with specific name. `environment` is null if not found.
output "environment" {
  value = one([
    for each in data.bastionzero_environments.example.environments
    : each if each.name == "example-env"
  ])
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `environments` (Attributes List) List of environments. (see [below for nested schema](#nestedatt--environments))
- `id` (String, Deprecated) Deprecated. Do not depend on this attribute. This attribute will be removed in the future.

<a id="nestedatt--environments"></a>
### Nested Schema for `environments`

Read-Only:

- `description` (String) The environment's description.
- `id` (String) The environment's unique ID.
- `is_default` (Boolean) If `true`, this environment is the default environment; `false` otherwise.
- `name` (String) The environment's name.
- `offline_cleanup_timeout_hours` (Number) The amount of time (in hours) to wait until offline targets are automatically removed by BastionZero (Defaults to `2160` hours [90 days]). If this value is `0`, then offline target cleanup is disabled.
- `organization_id` (String) The environment's organization's ID.
- `targets` (Attributes Map) Map of targets that belong to this environment. The map is keyed by a target's unique ID. (see [below for nested schema](#nestedatt--environments--targets))
- `time_created` (String) The time this environment was created in BastionZero formatted as a UTC timestamp string in [RFC 3339](https://datatracker.ietf.org/doc/html/rfc3339) format.

<a id="nestedatt--environments--targets"></a>
### Nested Schema for `environments.targets`

Read-Only:

- `id` (String) The target's unique ID.
- `type` (String) The target's type (one of `Bzero`, `Cluster`, `DynamicAccessConfig`, `Web`, or `Db`).