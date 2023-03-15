---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bastionzero_groups Data Source - terraform-provider-bastionzero"
subcategory: ""
description: |-
  Get a list of all groups in your BastionZero organization.
---

# bastionzero_groups (Data Source)

Get a list of all groups in your BastionZero organization.



<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `groups` (Attributes List) List of groups in your organization. (see [below for nested schema](#nestedatt--groups))

<a id="nestedatt--groups"></a>
### Nested Schema for `groups`

Read-Only:

- `id` (String) The group's unique ID, as specified by the Identity Provider in which it is configured.
- `name` (String) The group's name.

