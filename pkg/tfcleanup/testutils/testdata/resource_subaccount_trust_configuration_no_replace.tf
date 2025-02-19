resource "btp_subaccount_trust_configuration" "trust_1" {
  auto_create_shadow_users = false
  available_for_user_logon = true
  description              = "Identity Authentication tenant terraform.accounts.ondemand.com used for platform users"
  domain                   = "terraform.accounts.ondemand.com"
  identity_provider        = "terraform.accounts.ondemand.com"
  name                     = "terraform.accounts.ondemand.com (platform users)"
  status                   = "active"
  subaccount_id            = "1234567890"
}
