generate_hcl "_terramate_generated_main.tf" {
  condition = tm_contains(terramate.stack.tags, "btp")
  content {

    resource "random_uuid" "uuid" {}

    locals {
      subaccount_name      = "${tm_upper(terramate.stack.tags[1])} ${var.project_name} (DO NOT DELETE)"
      subaccount_subdomain = join("-", [lower(replace("${tm_upper(terramate.stack.tags[1])}-${var.project_name}", " ", "-")), random_uuid.uuid.result])
      service_name_prefix  = lower(replace("${tm_upper(terramate.stack.tags[1])}-${var.project_name}", " ", "-"))
      subaccount_cf_org    = local.subaccount_subdomain
    }

    resource "btp_directory" "project_directory" {
      name        = "Directory for ${local.subaccount_name}"
      description = "This is a directory with features."
      features    = ["DEFAULT", "ENTITLEMENTS", "AUTHORIZATIONS"]
    }

    resource "btp_directory_entitlement" "alert_notification_service_standard" {
      directory_id = btp_directory.project_directory.id
      service_name = "alert-notification"
      plan_name    = "standard"
    }

    resource "btp_directory_entitlement" "feature_flags_service_lite" {
      directory_id = btp_directory.project_directory.id
      service_name = "feature-flags"
      plan_name    = "lite"
    }

    resource "btp_directory_entitlement" "feature_flags_dashboard_app" {
      directory_id = btp_directory.project_directory.id
      service_name = "feature-flags-dashboard"
      plan_name    = "dashboard"
    }

    resource "btp_subaccount" "project_subaccount" {
      name         = local.subaccount_name
      subdomain    = local.subaccount_subdomain
      region       = var.subaccount_region
      parent_id    = btp_directory.project_directory.id
      beta_enabled = tm_ternary(terramate.stack.tags[1] == "dev", true, false)
      usage        = tm_ternary(terramate.stack.tags[1] == "prod", "USED_FOR_PRODUCTION", "NOT_USED_FOR_PRODUCTION")
      labels = {
        "stage"      = [tm_upper(terramate.stack.tags[1])]
        "costcenter" = [var.project_costcenter]
      }
      depends_on = [btp_directory_entitlement.alert_notification_service_standard, btp_directory_entitlement.feature_flags_service_lite, btp_directory_entitlement.feature_flags_dashboard_app]
    }

    resource "btp_subaccount_entitlement" "alert_notification_service_standard" {
      subaccount_id = btp_subaccount.project_subaccount.id
      service_name  = "alert-notification"
      plan_name     = "standard"
    }

    resource "btp_subaccount_entitlement" "feature_flags_service_lite" {
      subaccount_id = btp_subaccount.project_subaccount.id
      service_name  = "feature-flags"
      plan_name     = "lite"
    }

    resource "btp_subaccount_entitlement" "feature_flags_dashboard_app" {
      subaccount_id = btp_subaccount.project_subaccount.id
      service_name  = "feature-flags-dashboard"
      plan_name     = "dashboard"
    }

    data "btp_subaccount_service_plan" "alert_notification_service_standard" {
      subaccount_id = btp_subaccount.project_subaccount.id
      name          = "standard"
      offering_name = "alert-notification"
      depends_on    = [btp_subaccount_entitlement.alert_notification_service_standard]
    }

    resource "btp_subaccount_service_instance" "alert_notification_service_standard" {
      subaccount_id  = btp_subaccount.project_subaccount.id
      serviceplan_id = data.btp_subaccount_service_plan.alert_notification_service_standard.id
      name           = "${local.service_name_prefix}-alert-notification"
    }

    resource "btp_subaccount_subscription" "feature_flags_dashboard_app" {
      subaccount_id = btp_subaccount.project_subaccount.id
      app_name      = "feature-flags-dashboard"
      plan_name     = "dashboard"
      depends_on    = [btp_subaccount_entitlement.feature_flags_dashboard_app]
    }

    data "btp_subaccount_environments" "all" {
      subaccount_id = btp_subaccount.project_subaccount.id
    }
    resource "terraform_data" "cf_landscape_label" {
      input = length(var.cf_landscape_label) > 0 ? var.cf_landscape_label : [for env in data.btp_subaccount_environments.all.values : env if env.service_name == "cloudfoundry" && env.environment_type == "cloudfoundry"][0].landscape_label
    }

    resource "btp_subaccount_environment_instance" "cloudfoundry" {
      subaccount_id    = btp_subaccount.project_subaccount.id
      name             = local.subaccount_cf_org
      environment_type = "cloudfoundry"
      service_name     = "cloudfoundry"
      plan_name        = "standard"
      landscape_label  = terraform_data.cf_landscape_label.output
      parameters = jsonencode({
        instance_name = local.subaccount_cf_org
      })
    }

  }
}
