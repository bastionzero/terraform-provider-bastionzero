terraform {
  required_providers {
    bastionzero = {
      source  = "bastionzero/bastionzero"
      version = "~> 0.0"
    }
  }
}

# An empty provider configuration assumes the BASTIONZERO_API_SECRET environment
# variable is set. The provider uses the environment variable's value as the
# `api_secret` field.
provider "bastionzero" {}