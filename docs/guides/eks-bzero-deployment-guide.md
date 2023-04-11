---
page_title: "AWS EKS Deployment Guide"
subcategory: ""
description: |-
This guide explains how to deploy the BastionZero bzero agent on an AWS EKS cluster and create the policy required in order to connect
---

# AWS EKS Deployment Guide

This guide explains how to deploy the BastionZero
[`bzero`](https://github.com/bastionzero/bzero) agent on an AWS EKS cluster and
create the policy required in order to connect. By the end of this guide, your
cluster will be
[autodiscovered](https://docs.bastionzero.com/docs/deployment/installing-the-agent#autodiscovery)
as a cluster target by BastionZero, and you will be able to connect to it using
the [`zli`](https://github.com/bastionzero/zli).

This document is split into modular sections to help guide you through the
entire Terraform configuration. 

## Before you begin

* You must be an administrator of your BastionZero organization in order to
register a target.
* Create a BastionZero [API
key](https://docs.bastionzero.com/docs/admin-guide/authorization#creating-an-api-key)
in order to configure the Terraform provider. Manage your API keys at the API
key panel found [here](https://cloud.bastionzero.com/admin/apikeys).
* Create a BastionZero [registration API
key](https://docs.bastionzero.com/docs/admin-guide/authorization#registration-api-keys)
in order to register the EKS cluster as a target in your BastionZero
organization.
* Use Terraform 1.x or higher.
* You must have an EKS cluster that is already deployed and ready to use.

This guide assumes you have basic knowledge of the AWS Terraform provider, Helm
Terraform provider, and Kubernetes Terraform provider. All three providers are
used to install the Bzero agent in your EKS cluster. Please note that
BastionZero supports other cloud providers as well. Use this guide as a model to
register your Kubernetes clusters at other cloud providers.

## Setup

First, we setup the BastionZero Terraform provider, the [AWS Terraform
provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#authentication-and-configuration),
the [Helm Terraform
provider](https://registry.terraform.io/providers/hashicorp/helm/latest/docs#authentication),
and the [Kubernetes Terraform
provider](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs#authentication).

```terraform
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
    helm = {
      source  = "hashicorp/helm"
      version = "2.9.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.19.0"
    }
  }
}

# Configure the AWS provider.
provider "aws" {
  region = "us-east-1"
}

variable "cluster_id" {
  type        = string
  description = "AWS EKS cluster ID"
  nullable    = false
}

data "aws_eks_cluster" "cluster" {
  name = var.cluster_id
}

data "aws_eks_cluster_auth" "cluster" {
  name = var.cluster_id
}

# Configure the Kubernetes provider
provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  token                  = data.aws_eks_cluster_auth.cluster.token
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
}

# Configure the Helm provider
provider "helm" {
  kubernetes {
    host                   = data.aws_eks_cluster.cluster.endpoint
    token                  = data.aws_eks_cluster_auth.cluster.token
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  }
}

# Configure the BastionZero provider. An empty provider configuration assumes
# the BASTIONZERO_API_SECRET environment variable is set. The provider uses the
# environment variable's value as the `api_secret` field.
provider "bastionzero" {}
```

-> **Note** Version constraints are included for the sake of completeness.
Please change the version constraints as you see fit. Keep in mind though this
guide is written assuming the the versions defined in the Terraform snippet
above.

For security purposes, we choose to configure the BastionZero provider using an
environment variable as it's more secure than hardcoding the secret in the
Terraform file itself.

Set the `BASTIONZERO_API_SECRET` environment variable to the API key's secret
that you created in the [previous step](#before-you-begin) before running
`terraform apply`.

```sh
export BASTIONZERO_API_SECRET=api-secret
```

In addition, please set `var.cluster_id` to the ID of your EKS cluster. An easy
way to supply this value is through an environment variable.

```sh
export TF_VAR_cluster_id=cluster-id
```

Otherwise, use your preferred
[method](https://developer.hashicorp.com/terraform/language/values/variables#assigning-values-to-root-module-variables)
for supplying a value to an input variable.

-> **Note** You can find your cluster ID using the
[`aws`](https://aws.amazon.com/cli/) CLI. Run `aws eks list-clusters` to list
the cluster IDs available in your AWS account.

## Create an environment

All BastionZero targets belong to a single environment. Environments help
organize a collection of targets. They're especially useful when creating
BastionZero policy when you wish to apply the same set of policy access rules to
a group of targets.

Let's create a [`bastionzero_environment`](../resources/environment) to contain
our EKS cluster once it is autodiscovered and registered as a cluster target.

```terraform
resource "bastionzero_environment" "env" {
  name        = "example-env"
  description = "Environment managed by Terraform."
}
```

## Install the bzero agent via Helm

Next, we'll use the Helm provider and the [Bctl
quickstart](https://github.com/bastionzero/charts) chart to install the Bzero
agent in your Kubernetes cluster. We'll also use the Kubernetes provider to
create a separate namespace to contain the Bzero agent and associated Kubernetes
resources; this is best practice, and we recommend you always install the Bzero
agent in a separate namespace.

```terraform
variable "bzero_reg_secret" {
  type        = string
  description = "BastionZero registration secret used to register a target."
  sensitive   = true
  nullable    = false
}

variable "agent_version" {
  type        = string
  description = "Agent version to install."
  nullable    = false
  default     = "latest"
}

resource "kubernetes_namespace" "bastionzero_namespace" {
  metadata {
    name = "bastionzero"
  }
}

resource "kubernetes_secret" "registration_key_secret" {
  metadata {
    name      = "bctl-agent-registration-key"
    namespace = kubernetes_namespace.bastionzero_namespace.metadata[0].name
  }

  data = {
    "api-key" = var.bzero_reg_secret
  }
}

resource "helm_release" "bctl_agent" {
  name       = "bctl-agent"
  repository = "https://bastionzero.github.io/charts/"
  chart      = "bctl-quickstart"
  # version = Set this if you want to install a specifc version of the chart.
  # Otherwise, the latest chart is used.

  set {
    name  = "clusterName"
    value = var.cluster_id
  }

  set {
    name  = "apiKeyExistingSecret"
    value = kubernetes_secret.registration_key_secret.metadata[0].name
  }

  set {
    name  = "image.agentImageTag"
    value = var.agent_version
  }

  set {
    name  = "environmentId"
    value = bastionzero_environment.env.id
  }

  reuse_values = true

  namespace = kubernetes_namespace.bastionzero_namespace.metadata[0].name
}
```

Optionally, configure the `agent_version` input variable using your preferred
[method](https://developer.hashicorp.com/terraform/language/values/variables#assigning-values-to-root-module-variables)
for supplying a value to an input variable. Otherwise, accept the default value
which installs the latest agent version. Keep up to date with Bzero agent
updates [here](https://github.com/bastionzero/bzero/releases).

Configure the `bzero_reg_secret` input variable by setting an [environment
variable](https://developer.hashicorp.com/terraform/language/values/variables#environment-variables)
before running `terraform apply`.

```sh
export TF_VAR_bzero_reg_secret=reg-secret
```

~> **Warning** The registration secret is sensitive data. If a malicious
attacker obtains this credential, they could register their own instances as
targets in your BastionZero organization. Once the registration secret is used
in a Terraform module (e.g. fetched via a data source), it is stored in the
Terraform state file. Please protect your state files accordingly. See
HashiCorp's article about managing sensitive data in Terraform state
[here](https://developer.hashicorp.com/terraform/language/state/sensitive-data).

## Control access via policy

Your EKS cluster should appear in the
["Targets"](https://cloud.bastionzero.com/admin/targets) list shortly. Now let's
create a Kubernetes policy to control who in your organization has access to
your target.

Let's also create an example Kubernetes cluster role binding to the built-in
[`view`
ClusterRole](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles),
and refer to the Kubernetes RBAC username in the BastionZero policy that we'll
create. Feel free to omit this step and use existing usernames (that have
associated role bindings) that you've already configured via Kubernetes RBAC in
your cluster.

```terraform
# Create cluster rolebinding to the built-in `view` ClusterRole
resource "kubernetes_cluster_role_binding" "viewer_cluster_role_binding" {
  metadata {
    name = "viewer-cluster-rolebinding"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    # A default cluster role that comes built-in with Kubernetes installations.
    # It allows read-access to non-sensitive information. See
    # https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles
    # for a full description.
    name = "view"
  }
  subject {
    kind = "User"
    # The cluster username that should be specified in BastionZero Kubernetes
    # policy (`cluster_users`)
    name      = "viewer"
    api_group = "rbac.authorization.k8s.io"
  }
}

# Get all users and groups
data "bastionzero_users" "u" {}
data "bastionzero_groups" "g" {}

locals {
  # Define, by email address, users to add to the policy
  users = ["alice@example.com", "bob@example.com", "charlie@example.com"]
  # Define, by name, the groups to add to the policy
  groups = ["Engineering"]
}

resource "bastionzero_kubernetes_policy" "example" {
  name        = "example-policy"
  description = "Policy managed by Terraform."
  subjects = [
    for each in data.bastionzero_users.u.users
    : { id = each.id, type = each.type } if contains(local.users, each.email)
  ]
  groups = [
    for each in data.bastionzero_groups.g.groups
    : { id = each.id, name = each.name } if contains(local.groups, each.name)
  ]
  # Apply this policy to the environment created earlier
  environments = [bastionzero_environment.env.id]

  # Permit access as the "viewer" Kubernetes RBAC user, and therefore assume all
  # the permissions granted to the "viewer" user as defined by in-cluster
  # RoleBindings and ClusterRoleBindings.
  #
  # We created a ClusterRoleBinding for this user using the Kubernetes provider
  # above.
  cluster_users = ["viewer"]
}
```

See the [`bastionzero_kubernetes_policy`](../resources/kubernetes_policy)
resource documentation to learn about all the available options.

Specify the users' emails in your organization to whom you wish to give access
to in `local.users`.

If you've [synced groups](https://cloud.bastionzero.com/admin/integrations) from
your IdP (see the [SSO
Management](https://docs.bastionzero.com/docs/admin-guide/authentication/sso-management)
guide for more information), then specify the groups in your organization to
whom you wish to give access to in `local.groups`.

## Next steps

Congratulations! You've successfully registered an EKS cluster to your
BastionZero organization, and you've restricted access to it via BastionZero
policy.

If you experienced any issues with the BastionZero Terraform provider, please
open a ticket at
https://github.com/bastionzero/terraform-provider-bastionzero/issues. For
assistance on any issue with using or installing BastionZero, please reach out
to support@bastionzero.com.

Here are some suggestions on what to try next.

### Connect to your target using the `zli`

Download and install the `zli` to connect to your target. Follow the
instructions [here](https://github.com/bastionzero/zli#install).

```sh
zli login
zli connect viewer@cluster-id
kubectl get pods -A
```

Replace `cluster-id` with your EKS cluster ID.

### Setup a JIT policy

Enable the [Slack
integration](https://docs.bastionzero.com/docs/automation-and-integrations/slack)
in your BastionZero organization, so that you can write just-in-time (JIT)
policies to grant temporary access to your BastionZero target subject to
administrator approval.

Use the [`bastionzero_jit_policy`](../resources/jit_policy) resource to manage
the policy in Terraform. Learn more about JIT policies
[here](https://docs.bastionzero.com/docs/admin-guide/authorization#just-in-time).