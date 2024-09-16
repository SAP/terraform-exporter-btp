
###
# Resource: SUBACCOUNT
###
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
resource "btp_subaccount" "btptfexporter-validate" {
  beta_enabled = false
  description  = null
  labels       = null
  name         = "btptfexporter-validate"
  parent_id    = "cb0997e4-a8e3-4c4c-9ed2-a90e3fff8c3a"
  region       = "us10"
  subdomain    = "btptfexporter-validate-xhyu0cme"
  usage        = "NOT_USED_FOR_PRODUCTION"
}

###
# Resource: SUBACCOUNT ENTITLEMENTS
###
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,feature-flags-dashboard,dashboard"
resource "btp_subaccount_entitlement" "feature-flags-dashboard_dashboard" {
  amount        = 1
  plan_name     = "dashboard"
  service_name  = "feature-flags-dashboard"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,service-manager,service-operator-access"
resource "btp_subaccount_entitlement" "service-manager_service-operator-access" {
  amount        = 1
  plan_name     = "service-operator-access"
  service_name  = "service-manager"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,xsuaa,space"
resource "btp_subaccount_entitlement" "xsuaa_space" {
  amount        = 1
  plan_name     = "space"
  service_name  = "xsuaa"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,service-manager,subaccount-admin"
resource "btp_subaccount_entitlement" "service-manager_subaccount-admin" {
  amount        = 1
  plan_name     = "subaccount-admin"
  service_name  = "service-manager"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,xsuaa,broker"
resource "btp_subaccount_entitlement" "xsuaa_broker" {
  amount        = 1
  plan_name     = "broker"
  service_name  = "xsuaa"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,application-logs,lite"
resource "btp_subaccount_entitlement" "application-logs_lite" {
  amount        = 1
  plan_name     = "lite"
  service_name  = "application-logs"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,xsuaa,application"
resource "btp_subaccount_entitlement" "xsuaa_application" {
  amount        = 1
  plan_name     = "application"
  service_name  = "xsuaa"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,print,receiver"
resource "btp_subaccount_entitlement" "print_receiver" {
  amount        = 1
  plan_name     = "receiver"
  service_name  = "print"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,html5-apps-repo,app-host"
resource "btp_subaccount_entitlement" "html5-apps-repo_app-host" {
  amount        = 1
  plan_name     = "app-host"
  service_name  = "html5-apps-repo"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,credstore,proxy"
resource "btp_subaccount_entitlement" "credstore_proxy" {
  amount        = 1
  plan_name     = "proxy"
  service_name  = "credstore"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,connectivity,lite"
resource "btp_subaccount_entitlement" "connectivity_lite" {
  amount        = 1
  plan_name     = "lite"
  service_name  = "connectivity"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,saas-registry,application"
resource "btp_subaccount_entitlement" "saas-registry_application" {
  amount        = 1
  plan_name     = "application"
  service_name  = "saas-registry"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,mdo-one-mds-master,standard"
resource "btp_subaccount_entitlement" "mdo-one-mds-master_standard" {
  amount        = 1
  plan_name     = "standard"
  service_name  = "mdo-one-mds-master"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,content-agent,application"
resource "btp_subaccount_entitlement" "content-agent_application" {
  amount        = 1
  plan_name     = "application"
  service_name  = "content-agent"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,sap-identity-services-onboarding,default"
resource "btp_subaccount_entitlement" "sap-identity-services-onboarding_default" {
  amount        = 1
  plan_name     = "default"
  service_name  = "sap-identity-services-onboarding"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,feature-flags,lite"
resource "btp_subaccount_entitlement" "feature-flags_lite" {
  amount        = 1
  plan_name     = "lite"
  service_name  = "feature-flags"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,one-mds-master,standard"
resource "btp_subaccount_entitlement" "one-mds-master_standard" {
  amount        = 1
  plan_name     = "standard"
  service_name  = "one-mds-master"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,destination,lite"
resource "btp_subaccount_entitlement" "destination_lite" {
  amount        = 1
  plan_name     = "lite"
  service_name  = "destination"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,auditlog-api,default"
resource "btp_subaccount_entitlement" "auditlog-api_default" {
  amount        = 1
  plan_name     = "default"
  service_name  = "auditlog-api"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,service-manager,subaccount-audit"
resource "btp_subaccount_entitlement" "service-manager_subaccount-audit" {
  amount        = 1
  plan_name     = "subaccount-audit"
  service_name  = "service-manager"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,auditlog-management,default"
resource "btp_subaccount_entitlement" "auditlog-management_default" {
  amount        = 1
  plan_name     = "default"
  service_name  = "auditlog-management"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,one-mds,sap-integration"
resource "btp_subaccount_entitlement" "one-mds_sap-integration" {
  amount        = 1
  plan_name     = "sap-integration"
  service_name  = "one-mds"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,cias,standard"
resource "btp_subaccount_entitlement" "cias_standard" {
  amount        = 1
  plan_name     = "standard"
  service_name  = "cias"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,html5-apps-repo,app-runtime"
resource "btp_subaccount_entitlement" "html5-apps-repo_app-runtime" {
  amount        = 1
  plan_name     = "app-runtime"
  service_name  = "html5-apps-repo"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,cias,oauth2"
resource "btp_subaccount_entitlement" "cias_oauth2" {
  amount        = 1
  plan_name     = "oauth2"
  service_name  = "cias"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,content-agent-ui,free"
resource "btp_subaccount_entitlement" "content-agent-ui_free" {
  amount        = 1
  plan_name     = "free"
  service_name  = "content-agent-ui"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,xsuaa,apiaccess"
resource "btp_subaccount_entitlement" "xsuaa_apiaccess" {
  amount        = 1
  plan_name     = "apiaccess"
  service_name  = "xsuaa"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,identity,application"
resource "btp_subaccount_entitlement" "identity_application" {
  amount        = 1
  plan_name     = "application"
  service_name  = "identity"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,feature-flags,standard"
resource "btp_subaccount_entitlement" "feature-flags_standard" {
  amount        = 1
  plan_name     = "standard"
  service_name  = "feature-flags"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,service-manager,container"
resource "btp_subaccount_entitlement" "service-manager_container" {
  amount        = 1
  plan_name     = "container"
  service_name  = "service-manager"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,content-agent,standard"
resource "btp_subaccount_entitlement" "content-agent_standard" {
  amount        = 1
  plan_name     = "standard"
  service_name  = "content-agent"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,connectivity,connectivity_proxy"
resource "btp_subaccount_entitlement" "connectivity_connectivity_proxy" {
  amount        = 1
  plan_name     = "connectivity_proxy"
  service_name  = "connectivity"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,autoscaler,standard"
resource "btp_subaccount_entitlement" "autoscaler_standard" {
  amount        = 1
  plan_name     = "standard"
  service_name  = "autoscaler"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,service-manager,global-offerings-audit"
resource "btp_subaccount_entitlement" "service-manager_global-offerings-audit" {
  amount        = 2
  plan_name     = "global-offerings-audit"
  service_name  = "service-manager"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

###
# Resource: SUBACCOUNT ENVIRONMENT INSTANCES
###
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,8D363D57-86F1-4FF4-B3F8-C1FB7FFD34AC"
resource "btp_subaccount_environment_instance" "cloudfoundry" {
  environment_type = "cloudfoundry"
  landscape_label  = "cf-us10"
  name             = "btptfexporter-validate-xhyu0cme_btptfexporter-validate"
  parameters = jsonencode({
    instance_name = "terraform-integration-prod_btptfexporter-validate-xhyu0cme"
  })
  plan_name     = "standard"
  service_name  = "cloudfoundry"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
  timeouts      = null
}

###
# Resource: SUBACCOUNT  SUBSCRIPTIONS
###
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,feature-flags-dashboard,dashboard"
resource "btp_subaccount_subscription" "feature_flags_dashboard" {
  app_name      = "feature-flags-dashboard"
  parameters    = jsonencode({})
  plan_name     = "dashboard"
  subaccount_id = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
  timeouts      = null
}

###
# Resource: SUBACCOUNT TRUST CONFIGURATIONS
###
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,sap.default"
resource "btp_subaccount_trust_configuration" "trust0" {
  auto_create_shadow_users = false
  available_for_user_logon = true
  description              = null
  domain                   = null
  identity_provider        = ""
  link_text                = "Default Identity Provider"
  name                     = "sap.default"
  status                   = "active"
  subaccount_id            = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}

# __generated__ by Terraform from "cf18dcef-5b15-4638-bd98-73f10c7f4f3a,terraformeds-platform"
resource "btp_subaccount_trust_configuration" "trust1" {
  auto_create_shadow_users = false
  available_for_user_logon = false
  description              = "Identity Authentication tenant terraformeds.accounts.ondemand.com used for platform users"
  domain                   = "terraformeds.accounts.ondemand.com"
  identity_provider        = "terraformeds.accounts.ondemand.com"
  link_text                = null
  name                     = "terraformeds.accounts.ondemand.com (platform users)"
  status                   = "active"
  subaccount_id            = "cf18dcef-5b15-4638-bd98-73f10c7f4f3a"
}
