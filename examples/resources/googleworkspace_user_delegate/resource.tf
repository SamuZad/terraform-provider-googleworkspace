# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0
resource "googleworkspace_user" "manager" {
  primary_email = "manager@example.com"
  password      = "34819d7beeabb9260a5c854bc85b3e44"
  hash_function = "MD5"

  name {
    family_name = "Scott"
    given_name  = "Michael"
  }
}

resource "googleworkspace_user" "assistant" {
  primary_email = "assistant@example.com"
  password      = "34819d7beeabb9260a5c854bc85b3e44"
  hash_function = "MD5"

  name {
    family_name = "Martin"
    given_name  = "Pam"
  }
}

# Grant the assistant delegated access to the manager's Gmail mailbox.
resource "googleworkspace_user_delegate" "assistant_to_manager" {
  user_id        = googleworkspace_user.manager.primary_email
  delegate_email = googleworkspace_user.assistant.primary_email
}
