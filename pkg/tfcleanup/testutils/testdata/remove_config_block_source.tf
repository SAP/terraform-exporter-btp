resource "btp_subaccount" "subaccount_0" {
  beta_enabled = true
  labels = {
    costcenter = ["12345"]
    stage      = ["DEV"]
  }
  name      = "DEV Inttest Exporter (DO NOT DELETE)"
  region    = "us10"
  subdomain = "dev-inttest-exporter-8b2967c9-c3f3-68a4-4604-16702098e1d1"
  usage     = "NOT_USED_FOR_PRODUCTION"
}
resource "btp_subaccount_trust_configuration" "trust_0" {
  auto_create_shadow_users = false
  available_for_user_logon = true
  link_text                = "Default Identity Provider"
  name                     = "sap.default"
  status                   = "active"
  subaccount_id            = "1234567890"
}
