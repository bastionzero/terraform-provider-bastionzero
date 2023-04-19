data "bastionzero_user" "example_by_id" {
  id = "<user-id>"
}

data "bastionzero_user" "example_by_email" {
  id = "alice@example.com"
}