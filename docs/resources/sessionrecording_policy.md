---
page_title: "bastionzero_sessionrecording_policy Resource - terraform-provider-bastionzero"
subcategory: "policy"
description: |-
  Provides a BastionZero session recording policy. Session recording policies govern whether users' I/O during shell connections are recorded.
---

# bastionzero_sessionrecording_policy (Resource)

Provides a BastionZero session recording policy. Session recording policies govern whether users' I/O during shell connections are recorded.

Learn more about session recording policies [here](https://docs.bastionzero.com/docs/admin-guide/authorization#session-recording).

~> **Note on policy name** All policies (of any type) must have a unique name. If the
configured [`name`](#name) is not unique, an error is thrown.

## Example Usage

Enable session recording for all users and service accounts:

```terraform
# Get all users and service accounts
data "bastionzero_users" "u" {}
data "bastionzero_service_accounts" "u" {}

resource "bastionzero_sessionrecording_policy" "example" {
  name        = "example-policy"
  description = "Policy managed by Terraform."
  subjects = [
    for each in concat(
      data.bastionzero_users.u.users,
      data.bastionzero_service_accounts.u.service_accounts
    )
    : { id = each.id, type = each.type }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The policy's name.

### Optional

- `description` (String) The policy's description.
- `groups` (Attributes Set) Set of Identity Provider (IdP) groups that this policy applies to. (see [below for nested schema](#nestedatt--groups))
- `record_input` (Boolean) If true, then in addition to session output, session input should be recorded. If false, then only session output should be recorded (Defaults to false).
- `subjects` (Attributes Set) Set of subjects that this policy applies to. (see [below for nested schema](#nestedatt--subjects))

### Read-Only

- `id` (String) The policy's unique ID.
- `type` (String) The policy's type (constant value `SessionRecording`).

<a id="nestedatt--groups"></a>
### Nested Schema for `groups`

Required:

- `id` (String) The group's unique ID.
- `name` (String) The group's name.


<a id="nestedatt--subjects"></a>
### Nested Schema for `subjects`

Required:

- `id` (String) The subject's unique ID.
- `type` (String) The subject's type (one of `User`, `ApiKey`, or `ServiceAccount`).

## Import

Import is supported using the following syntax:

```shell
# Policy can be imported by specifying the unique identifier.
terraform import bastionzero_sessionrecording_policy.example "feb58c89-35f3-4615-9300-ef64aa944c8b"
```