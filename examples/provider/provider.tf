terraform {
  required_providers {
    bastionzero = {
      source  = "bastionzero/bastionzero"
      version = "~> 0.0"
    }
  }
}

# Set the variable value in *.tfvars file
# or using -var="bz_api_secret=..." CLI option
variable "bz_api_secret" {}

# Configure the BastionZero Provider
provider "bastionzero" {
  api_secret = var.bz_api_secret
}