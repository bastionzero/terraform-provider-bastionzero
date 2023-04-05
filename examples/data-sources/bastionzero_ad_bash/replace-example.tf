# Using Terraform >= 1.x syntax

data "bastionzero_ad_bash" "example" {
  environment_id     = "<environment-id>"
  target_name_option = "BashHostName"
}

locals {
  # This is only an example. We recommend to fetch this secret from your
  # preferred secrets manager.
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
