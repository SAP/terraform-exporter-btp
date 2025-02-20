resource "btp_subaccount_role" "role_6" {
  app_id             = "destination-xsappname!b62"
  description        = "Manage destination configurations, certificates and signing keys for SAML assertions issued by the Destination service on instance Level"
  name               = "Destination Administrator Instance"
  role_template_name = "Destination_Administrator_Instance"
  subaccount_id      = "1234567890"
}
