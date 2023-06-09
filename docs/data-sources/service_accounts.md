---
page_title: "bastionzero_service_accounts Data Source - terraform-provider-bastionzero"
subcategory: "Service Account"
description: |-
  Get a list of all service accounts in your BastionZero organization. A service account is a Google, Azure, or generic service account that integrates with BastionZero by sharing its JSON Web Key Set (JWKS) URL. The headless authentication closely follows the OpenID Connect (OIDC) protocol.
---

# bastionzero_service_accounts (Data Source)

Get a list of all service accounts in your BastionZero organization. A service account is a Google, Azure, or generic service account that integrates with BastionZero by sharing its JSON Web Key Set (JWKS) URL. The headless authentication closely follows the OpenID Connect (OIDC) protocol.

See the [Service Accounts
Management](https://docs.bastionzero.com/docs/admin-guide/authentication/service-accounts-management)
guide to learn how to configure service accounts with BastionZero.

-> **Note** You can use the [`bastionzero_service_account`](service_account)
data source to obtain metadata about a single service account if you already
know the `id`.

## Example Usage

```terraform
data "bastionzero_service_accounts" "example" {}

# Find all service accounts whose JWKS URLs contain "google"
output "google_sas" {
  value = [
    for each in data.bastionzero_service_accounts.example.service_accounts
    : each if can(regex("google", each.jwks_url))
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String, Deprecated) Deprecated. Do not depend on this attribute. This attribute will be removed in the future.
- `service_accounts` (Attributes List) List of service accounts. (see [below for nested schema](#nestedatt--service_accounts))

<a id="nestedatt--service_accounts"></a>
### Nested Schema for `service_accounts`

Read-Only:

- `created_by` (String) Unique identifier for the subject that created this service account.
- `email` (String) The service account's email address.
- `enabled` (Boolean) If `true`, the service account is currently enabled; `false` otherwise.
- `external_id` (String) The service account's unique per service provider identifier provided by the user during creation.
- `id` (String) The service account's unique ID.
- `is_admin` (Boolean) If `true`, the service account is an administrator; `false` otherwise.
- `jwks_url` (String) The service account's publicly available JWKS URL that provides the public key that can be used to verify the tokens signed by the private key of this service account.
- `jwks_url_pattern` (String) A URL pattern that all service accounts of the same service account provider follow in their JWKS URL.
- `last_login` (String) The time this service account last logged into BastionZero formatted as a UTC timestamp string in [RFC 3339](https://datatracker.ietf.org/doc/html/rfc3339) format. Null if the service account has never logged in.
- `organization_id` (String) The service account's organization's ID.
- `time_created` (String) The time this service account was created in BastionZero formatted as a UTC timestamp string in [RFC 3339](https://datatracker.ietf.org/doc/html/rfc3339) format.
- `type` (String) The subject's type (constant value `ServiceAccount`).