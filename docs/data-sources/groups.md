---
page_title: "bastionzero_groups Data Source - terraform-provider-bastionzero"
subcategory: "Group"
description: |-
  Get a list of all groups in your BastionZero organization. A group is an Identity provider (IdP) group synced to BastionZero.
---

# bastionzero_groups (Data Source)

Get a list of all groups in your BastionZero organization. A group is an Identity provider (IdP) group synced to BastionZero.

This data source is useful when creating policies so that the policy can apply
to a dynamic set of users depending on the user's group membership.

Syncing groups from your IdP is configured on the [App
Integrations](https://cloud.bastionzero.com/admin/integrations) page. See the
[SSO
Management](https://docs.bastionzero.com/docs/admin-guide/authentication/sso-management)
guide for more information.

## Example Usage

```terraform
data "bastionzero_groups" "example" {}

# Find all groups whose names are equal to "Engineering"
output "engineering_groups" {
  value = [
    for each in data.bastionzero_groups.example.groups
    : each if each.name == "Engineering"
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `groups` (Attributes List) List of groups. (see [below for nested schema](#nestedatt--groups))
- `id` (String, Deprecated) Deprecated. Do not depend on this attribute. This attribute will be removed in the future.

<a id="nestedatt--groups"></a>
### Nested Schema for `groups`

Read-Only:

- `id` (String) The group's unique ID, as specified by the Identity Provider in which it is configured.
- `name` (String) The group's name.