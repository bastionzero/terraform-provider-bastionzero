---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bastionzero_proxy_policies Data Source - terraform-provider-bastionzero"
subcategory: ""
description: |-
  Get a list of all proxy policies in your BastionZero organization.
---

# bastionzero_proxy_policies (Data Source)

Get a list of all proxy policies in your BastionZero organization.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter_groups` (Set of String) Filters the list of policies to only those that contain the provided group ID(s).
- `filter_subjects` (Set of String) Filters the list of policies to only those that contain the provided subject ID(s).

### Read-Only

- `policies` (Attributes List) List of proxy policies. (see [below for nested schema](#nestedatt--policies))

<a id="nestedatt--policies"></a>
### Nested Schema for `policies`

Read-Only:

- `description` (String) The policy's description.
- `environments` (Set of String) Set of environments that this policy applies to.
- `groups` (Attributes Set) Set of Identity Provider (IdP) groups that this policy applies to. (see [below for nested schema](#nestedatt--policies--groups))
- `id` (String) The policy's unique ID.
- `name` (String) The policy's name.
- `subjects` (Attributes Set) Set of subjects that this policy applies to. (see [below for nested schema](#nestedatt--policies--subjects))
- `target_users` (Set of String) Set of Database usernames that this policy applies to. These usernames only affect policy decisions involving Db targets that have the SplitCert feature enabled.
- `targets` (Attributes Set) Set of targets that this policy applies to. (see [below for nested schema](#nestedatt--policies--targets))
- `type` (String) The policy's type (constant value "Proxy").

<a id="nestedatt--policies--groups"></a>
### Nested Schema for `policies.groups`

Read-Only:

- `id` (String) The group's unique ID.
- `name` (String) The group's name.


<a id="nestedatt--policies--subjects"></a>
### Nested Schema for `policies.subjects`

Read-Only:

- `id` (String) The subject's unique ID.
- `type` (String) The subject's type (one of "User", "ApiKey", or "ServiceAccount").


<a id="nestedatt--policies--targets"></a>
### Nested Schema for `policies.targets`

Read-Only:

- `id` (String) The target's unique ID.
- `type` (String) The target's type (one of "Db", or "Web").

