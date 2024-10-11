
###
# Resource: BTP_DIRECTORY_ROLE_COLLECTION
###
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform from "f366cc89-9fa7-42d8-a62e-84ccec2b16a1,Directory Viewer"
resource "btp_directory_role_collection" "directory_viewer" {
  description  = "Read-only access to the directory"
  directory_id = "f366cc89-9fa7-42d8-a62e-84ccec2b16a1"
  name         = "Directory Viewer"
  roles = [
    {
      name                 = "Directory Usage Reporting Viewer"
      role_template_app_id = "uas!b36585"
      role_template_name   = "Directory_Usage_Reporting_Viewer"
    },
    {
      name                 = "Directory Viewer"
      role_template_app_id = "cis-central!b14"
      role_template_name   = "Directory_Viewer"
    },
    {
      name                 = "User and Role Auditor"
      role_template_app_id = "xsuaa!t1"
      role_template_name   = "xsuaa_auditor"
    },
  ]
}

# __generated__ by Terraform from "f366cc89-9fa7-42d8-a62e-84ccec2b16a1,Directory Administrator"
resource "btp_directory_role_collection" "directory_administrator" {
  description  = "Administrative access to the directory"
  directory_id = "f366cc89-9fa7-42d8-a62e-84ccec2b16a1"
  name         = "Directory Administrator"
  roles = [
    {
      name                 = "Directory Admin"
      role_template_app_id = "cis-central!b14"
      role_template_name   = "Directory_Admin"
    },
    {
      name                 = "Directory Usage Reporting Viewer"
      role_template_app_id = "uas!b36585"
      role_template_name   = "Directory_Usage_Reporting_Viewer"
    },
    {
      name                 = "User and Role Administrator"
      role_template_app_id = "xsuaa!t1"
      role_template_name   = "xsuaa_admin"
    },
  ]
}
