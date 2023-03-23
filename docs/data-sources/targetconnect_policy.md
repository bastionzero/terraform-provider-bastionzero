---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bastionzero_targetconnect_policy Data Source - terraform-provider-bastionzero"
subcategory: ""
description: |-
  Get information on a BastionZero target connect policy.
---

# bastionzero_targetconnect_policy (Data Source)

Get information on a BastionZero target connect policy.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The policy's unique ID.

### Read-Only

- `description` (String) The policy's description.
- `environments` (Set of String) Set of environments that this policy applies to.
- `groups` (Attributes Set) Set of Identity Provider (IdP) groups that this policy applies to. (see [below for nested schema](#nestedatt--groups))
- `name` (String) The policy's name.
- `subjects` (Attributes Set) Set of subjects that this policy applies to. (see [below for nested schema](#nestedatt--subjects))
- `target_users` (Set of String) Set of Unix usernames that this policy applies to.
- `targets` (Attributes Set) Set of targets that this policy applies to. (see [below for nested schema](#nestedatt--targets))
- `type` (String) The policy's type (constant value "TargetConnect").
- `verbs` (Set of String) Set of actions allowed by this policy (one of "Shell", "FileTransfer", or "Tunnel").

<a id="nestedatt--groups"></a>
### Nested Schema for `groups`

Read-Only:

- `id` (String) The group's unique ID.
- `name` (String) The group's name.


<a id="nestedatt--subjects"></a>
### Nested Schema for `subjects`

Read-Only:

- `id` (String) The subject's unique ID.
- `type` (String) The subject's type (one of "User", "ApiKey", or "ServiceAccount").


<a id="nestedatt--targets"></a>
### Nested Schema for `targets`

Read-Only:

- `id` (String) The target's unique ID.
- `type` (String) The target's type (one of "Bzero", or "DynamicAccessConfig").

