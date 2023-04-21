---
page_title: "bastionzero_ad_bash Data Source - terraform-provider-bastionzero"
subcategory: ""
description: |-
  Get a bash script that can be used to install the latest production BastionZero agent (bzero https://github.com/bastionzero/bzero) on your targets.
---

# bastionzero_ad_bash (Data Source)

Get a bash script that can be used to install the latest production BastionZero agent ([`bzero`](https://github.com/bastionzero/bzero)) on your targets.

-> **Note** If you do not have a default global registration key selected at the
[API key panel](https://cloud.bastionzero.com/admin/apikeys), then the fetched
`script` does not contain the registration secret that is required to register
your targets with BastionZero. You must replace
`<REGISTRATION-SECRET-GOES-HERE>` with a valid [registration
secret](https://docs.bastionzero.com/docs/admin-guide/authorization#registration-api-keys)
before attempting to execute the script. This can be done by using the
[`replace`](https://www.terraform.io/language/functions/replace) function (see
example [below](#replace-example)).

## Example Usage

### Basic example

```terraform
data "bastionzero_ad_bash" "example" {
  environment_id     = "<environment-id>"
  target_name_option = "BashHostName"
}
```

### Register in the default environment

The following example fetches a script that registers your target in the default
environment.

```terraform
# Using Terraform >= 1.x syntax

# Get all environments
data "bastionzero_environments" "envs" {}

# Find the default environment which is guaranteed to exist
locals {
  default_env_id = one([
    for each in data.bastionzero_environments.envs.environments
    : each if each.is_default
  ]).id
}

data "bastionzero_ad_bash" "example" {
  environment_id     = local.default_env_id
  target_name_option = "BashHostName"
}
```

### Replace example

~> **Warning** The registration secret is sensitive data. If a malicious
attacker obtains this credential, they could register their own instances as
targets in your BastionZero organization. Once the registration secret is used
in a Terraform module (e.g. fetched via a data source), it is stored in the
Terraform state file. Please protect your state files accordingly. See
HashiCorp's article about managing sensitive data in Terraform state
[here](https://developer.hashicorp.com/terraform/language/state/sensitive-data).

```terraform
# Using Terraform >= 1.x syntax

data "bastionzero_ad_bash" "example" {
  environment_id     = "<environment-id>"
  target_name_option = "BashHostName"
}

locals {
  # This is only an example. We recommend to fetch this secret from your
  # preferred secrets manager. Do not expose a .tf file with your secret
  reg_key_secret = sensitive("<your-registration-key-secret>")

  # This script can be used during cloud-init (User data) when provisioning your
  # cloud instances. 
  script_ready_to_use = sensitive(
    replace(
      data.bastionzero_ad_bash.example.script,
      "<REGISTRATION-SECRET-GOES-HERE>",
      local.reg_key_secret
    )
  )
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The unique environment ID the target should associate with.
- `target_name_option` (String) The target name schema option to use during autodiscovery (one of `Timestamp`, `DigitalOceanMetadata`, `AwsEc2Metadata`, or `BashHostName`).

### Read-Only

- `id` (String, Deprecated) Deprecated. Do not depend on this attribute. This attribute will be removed in the future.
- `script` (String, Sensitive) Bash script that can be used to autodiscover a target.