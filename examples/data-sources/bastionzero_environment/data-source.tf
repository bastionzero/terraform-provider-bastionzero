data "bastionzero_environment" "example" {
  id = "7b2120a7-7bb1-4790-b924-334cc65cfc80"
}

output "example_env_targets" {
  value = data.bastionzero_environment.example.targets
}