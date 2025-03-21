// TERRAMATE: GENERATED AUTOMATICALLY DO NOT EDIT

resource "random_uuid" "uuid" {
}
locals {
  service_name_prefix = lower(replace("DEV-${var.project_name}", " ", "-"))
  subaccount_cf_org   = local.subaccount_subdomain
  subaccount_name     = "DEV ${var.project_name} (DO NOT DELETE)"
  subaccount_subdomain = join("-", [
    lower(replace("DEV-${var.project_name}", " ", "-")),
    random_uuid.uuid.result,
  ])
}
resource "btp_directory" "project_directory" {
  description = "This is a directory with features."
  features = [
    "DEFAULT",
    "ENTITLEMENTS",
    "AUTHORIZATIONS",
  ]
  name = "Directory for ${local.subaccount_name}"
}
resource "btp_directory_role_collection" "custom" {
  description  = "A custom Role Collection to validate filtering."
  directory_id = btp_directory.project_directory.id
  name         = "Directory Custom Role Collection"
  roles = [
    {
      name                 = "Directory Admin"
      role_template_app_id = "cis-central!b13"
      role_template_name   = "Directory_Admin"
    },
  ]
}
resource "btp_directory_role" "directory_viewer" {
  app_id             = "cis-central!b14"
  directory_id       = btp_directory.project_directory.id
  name               = "Custom Directory Viewer Role"
  role_template_name = "Directory_Viewer"
}
resource "btp_directory_entitlement" "alert_notification_service_standard" {
  directory_id = btp_directory.project_directory.id
  plan_name    = "standard"
  service_name = "alert-notification"
}
resource "btp_directory_entitlement" "feature_flags_service_lite" {
  directory_id = btp_directory.project_directory.id
  plan_name    = "lite"
  service_name = "feature-flags"
}
resource "btp_directory_entitlement" "feature_flags_dashboard_app" {
  directory_id = btp_directory.project_directory.id
  plan_name    = "dashboard"
  service_name = "feature-flags-dashboard"
}
resource "btp_subaccount" "project_subaccount" {
  beta_enabled = true
  depends_on = [
    btp_directory_entitlement.alert_notification_service_standard,
    btp_directory_entitlement.feature_flags_service_lite,
    btp_directory_entitlement.feature_flags_dashboard_app,
  ]
  labels = {
    "stage" = [
      "DEV",
    ]
    "costcenter" = [
      var.project_costcenter,
    ]
  }
  name      = local.subaccount_name
  parent_id = btp_directory.project_directory.id
  region    = var.subaccount_region
  subdomain = local.subaccount_subdomain
  usage     = "NOT_USED_FOR_PRODUCTION"
}
resource "btp_subaccount_role_collection" "custom" {
  description = "A custom Role Collection to validate filtering."
  name        = "Directory Custom Role Collection"
  roles = [
    {
      name                 = "Application Subscriptions Viewer"
      role_template_app_id = "cis-local!b4"
      role_template_name   = "Application_Subscriptions_Viewer"
    },
    {
      name                 = "Subaccount Viewer"
      role_template_app_id = "cis-local!b4"
      role_template_name   = "Subaccount_Viewer"
    },
    {
      name                 = "Subaccount_Service_Viewer"
      role_template_app_id = "service-manager!b1476"
      role_template_name   = "Subaccount_Service_Viewer"
    },
  ]
  subaccount_id = btp_subaccount.project_subaccount.id
}
resource "btp_subaccount_role" "xsuaa_auditor" {
  app_id             = "xsuaa!t8"
  name               = "Custom XSUAA Auditor Role"
  role_template_name = "xsuaa_auditor"
  subaccount_id      = btp_subaccount.project_subaccount.id
}
resource "btp_subaccount_entitlement" "alert_notification_service_standard" {
  plan_name     = "standard"
  service_name  = "alert-notification"
  subaccount_id = btp_subaccount.project_subaccount.id
}
resource "btp_subaccount_entitlement" "feature_flags_service_lite" {
  plan_name     = "lite"
  service_name  = "feature-flags"
  subaccount_id = btp_subaccount.project_subaccount.id
}
resource "btp_subaccount_entitlement" "feature_flags_dashboard_app" {
  plan_name     = "dashboard"
  service_name  = "feature-flags-dashboard"
  subaccount_id = btp_subaccount.project_subaccount.id
}
data "btp_subaccount_service_plan" "alert_notification_service_standard" {
  depends_on = [
    btp_subaccount_entitlement.alert_notification_service_standard,
  ]
  name          = "standard"
  offering_name = "alert-notification"
  subaccount_id = btp_subaccount.project_subaccount.id
}
resource "btp_subaccount_service_instance" "alert_notification_service_standard" {
  name           = "${local.service_name_prefix}-alert-notification"
  serviceplan_id = data.btp_subaccount_service_plan.alert_notification_service_standard.id
  subaccount_id  = btp_subaccount.project_subaccount.id
}
resource "btp_subaccount_service_binding" "my_binding" {
  name                = "${local.service_name_prefix}-sb-af"
  service_instance_id = btp_subaccount_service_instance.alert_notification_service_standard.id
  subaccount_id       = btp_subaccount.project_subaccount.id
}
resource "btp_subaccount_subscription" "feature_flags_dashboard_app" {
  app_name = "feature-flags-dashboard"
  depends_on = [
    btp_subaccount_entitlement.feature_flags_dashboard_app,
  ]
  plan_name     = "dashboard"
  subaccount_id = btp_subaccount.project_subaccount.id
}
data "btp_subaccount_environments" "all" {
  subaccount_id = btp_subaccount.project_subaccount.id
}
resource "terraform_data" "cf_landscape_label" {
  input = length(var.cf_landscape_label) > 0 ? var.cf_landscape_label : [for env in data.btp_subaccount_environments.all.values : env if env.service_name == "cloudfoundry" && env.environment_type == "cloudfoundry"][0].landscape_label
}
resource "btp_subaccount_environment_instance" "cloudfoundry" {
  environment_type = "cloudfoundry"
  landscape_label  = terraform_data.cf_landscape_label.output
  name             = local.subaccount_cf_org
  parameters = jsonencode({
    instance_name = local.subaccount_cf_org
  })
  plan_name     = "standard"
  service_name  = "cloudfoundry"
  subaccount_id = btp_subaccount.project_subaccount.id
}
