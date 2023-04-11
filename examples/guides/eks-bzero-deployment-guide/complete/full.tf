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
resource "bastionzero_environment" "env" {
  name        = "example-env"
  description = "Environment managed by Terraform."
}
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
} # Create cluster rolebinding to the built-in `view` ClusterRole
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
