---
page_title: "bastionzero_jit_policy Data Source - terraform-provider-bastionzero"
subcategory: "policy"
description: |-
  Get information on a BastionZero JIT policy. A JIT policy provides just in time access to targets.
---

# bastionzero_jit_policy (Data Source)

Get information on a BastionZero JIT policy. A JIT policy provides just in time access to targets.

This is useful if the policy in question is not managed by Terraform, or
you need to utilize any of the policy's data.

## Example Usage

Get the policy by ID:

```terraform
data "bastionzero_jit_policy" "example" {
  id = "<policy-id>"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The policy's unique ID.

### Read-Only

- `auto_approved` (Boolean) If true, then the policies created by this JIT policy will be automatically approved. If false, then policies will only be created based on request and approval from reviewers (Defaults to false).
- `child_policies` (Attributes Set) Set of policies that a JIT policy applies to. (see [below for nested schema](#nestedatt--child_policies))
- `description` (String) The policy's description.
- `duration` (Number) The amount of time (in minutes) after which the access granted by this JIT policy will expire (Defaults to 1 hour).
- `groups` (Attributes Set) Set of Identity Provider (IdP) groups that this policy applies to. (see [below for nested schema](#nestedatt--groups))
- `name` (String) The policy's name.
- `subjects` (Attributes Set) Set of subjects that this policy applies to. (see [below for nested schema](#nestedatt--subjects))
- `type` (String) The policy's type (constant value `JustInTime`).

<a id="nestedatt--child_policies"></a>
### Nested Schema for `child_policies`

Read-Only:

- `id` (String) The policy's unique ID.
- `name` (String) The policy's name.
- `type` (String) The policy's type (one of `TargetConnect`, `Kubernetes`, or `Proxy`).


<a id="nestedatt--groups"></a>
### Nested Schema for `groups`

Read-Only:

- `id` (String) The group's unique ID.
- `name` (String) The group's name.


<a id="nestedatt--subjects"></a>
### Nested Schema for `subjects`

Read-Only:

- `id` (String) The subject's unique ID.
- `type` (String) The subject's type (one of `User`, `ApiKey`, or `ServiceAccount`).