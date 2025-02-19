resource "btp_subaccount_trust_configuration" "trust_0" {
  auto_create_shadow_users = false
  available_for_user_logon = true
  description              = null
  domain                   = null
  identity_provider        = ""
  link_text                = "Default Identity Provider"
  name                     = "sap.default"
  status                   = "active"
  subaccount_id            = "1234567890"
}
