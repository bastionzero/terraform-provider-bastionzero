data "bastionzero_cluster_target" "example" {
  name = "example-cluster"
  timeouts = {
    # Wait up to 30 seconds to find target with name "example-cluster"
    read = "30s"
  }
}