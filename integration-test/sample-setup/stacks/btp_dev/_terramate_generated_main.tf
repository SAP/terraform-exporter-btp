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
