resource "btp_subaccount_role_collection" "rolecollection_3" {
  description = "Administrative access to the subaccount"
  name        = "Subaccount Administrator"
  roles = [
    {
      name                 = "Cloud Connector Administrator"
      role_template_app_id = "connectivity!b7"
      role_template_name   = "Cloud_Connector_Administrator"
    },
    {
      name                 = "Destination Administrator"
      role_template_app_id = "destination-xsappname!b62"
      role_template_name   = "Destination_Administrator"
    },
    {
      name                 = "Subaccount Service Administrator"
      role_template_app_id = "service-manager!b1476"
      role_template_name   = "Subaccount_Service_Administrator"
    },
    {
      name                 = "User and Role Administrator"
      role_template_app_id = "xsuaa!t8"
      role_template_name   = "xsuaa_admin"
    },
  ]
  subaccount_id = "123456"
  depends_on    = [btp_subaccount_role.role_2]
}
