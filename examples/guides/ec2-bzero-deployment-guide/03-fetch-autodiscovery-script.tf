# This is only an example. We recommend to fetch this secret from your preferred
# secrets manager. Do not expose a .tf file with your secret.
variable "bzero_reg_secret" {
  type        = string
  description = "BastionZero registration secret used to register a target."
  sensitive   = true
  nullable    = false
}

data "bastionzero_ad_bash" "ad_script" {
  environment_id     = bastionzero_environment.env.id
  target_name_option = "AwsEc2Metadata"
}

locals {
  ad_script = sensitive(
    replace(
      data.bastionzero_ad_bash.ad_script.script,
      "<REGISTRATION-SECRET-GOES-HERE>",
      var.bzero_reg_secret
    )
  )
}
