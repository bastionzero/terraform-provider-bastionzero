---
page_title: "bastionzero_environment Resource - terraform-provider-bastionzero"
subcategory: "Environment"
description: |-
  Provides a BastionZero environment. An environment is a collection of targets.
---

# bastionzero_environment (Resource)

Provides a BastionZero environment. An environment is a collection of targets.

~> **Note on offline target cleanup** An environment's
[`offline_cleanup_timeout_hours`](#offline_cleanup_timeout_hours) cannot exceed
4320 hours (180 days).

## Example Usage

Create an environment named `example-env`:

```terraform
resource "bastionzero_environment" "example" {
  name                          = "example-env"
  description                   = "Example environment"
  offline_cleanup_timeout_hours = 12
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The environment's name.

### Optional

- `description` (String) The environment's description.
- `offline_cleanup_timeout_hours` (Number) The amount of time (in hours) to wait until offline targets are automatically removed by BastionZero (Defaults to `2160` hours [90 days]). If this value is `0`, then offline target cleanup is disabled.

### Read-Only

- `id` (String) The environment's unique ID.
- `is_default` (Boolean) If `true`, this environment is the default environment; `false` otherwise.
- `organization_id` (String) The environment's organization's ID.
- `targets` (Attributes Map) Map of targets that belong to this environment. The map is keyed by a target's unique ID. (see [below for nested schema](#nestedatt--targets))
- `time_created` (String) The time this environment was created in BastionZero formatted as a UTC timestamp string in [RFC 3339](https://datatracker.ietf.org/doc/html/rfc3339) format.

<a id="nestedatt--targets"></a>
### Nested Schema for `targets`

Read-Only:

- `id` (String) The target's unique ID.
- `type` (String) The target's type (one of `Bzero`, `Cluster`, `DynamicAccessConfig`, `Web`, or `Db`).

## Import

Import is supported using the following syntax:

```shell
# Environment can be imported by specifying the unique identifier.
terraform import bastionzero_environment.example "01d7a020-2bbc-4ac0-b886-ac9e445d8ab1"
```