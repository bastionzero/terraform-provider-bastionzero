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