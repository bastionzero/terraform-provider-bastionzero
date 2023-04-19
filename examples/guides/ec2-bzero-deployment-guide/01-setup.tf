terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
    bastionzero = {
      source  = "bastionzero/bastionzero"
      version = "~> 0.0"
    }
  }
}

# Configure the AWS provider.
provider "aws" {
  region = "us-east-1"
}

# Configure the BastionZero provider. An empty provider configuration assumes
# the BASTIONZERO_API_SECRET environment variable is set. The provider uses the
# environment variable's value as the `api_secret` field.
provider "bastionzero" {}
