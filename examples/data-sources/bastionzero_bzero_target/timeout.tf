data "bastionzero_bzero_target" "example" {
  name = "example-target"
  timeouts = {
    # Wait up to 30 seconds to find target with name "example-target"
    read = "30s"
  }
}