resource "btp_subaccount_entitlement" "entitlement_3" {
  plan_name              = "api_access"
  service_name           = "connectivity"
  subaccount_id          = "b9b11ec2-40e3-454f-bd23-f5ec2fe4bac8"
  plan_unique_identifier = "connectivity-api_access"
}

resource "btp_subaccount_entitlement" "entitlement_4" {
  plan_name              = "standard"
  service_name           = "cloudfoundry"
  subaccount_id          = "b9b11ec2-40e3-454f-bd23-f5ec2fe4bac8"
  plan_unique_identifier = "cloudfoundry-standard"
}

resource "btp_subaccount_entitlement" "entitlement_0" {
  amount                 = 2
  plan_name              = "MEMORY"
  service_name           = "APPLICATION_RUNTIME"
  subaccount_id          = "b9b11ec2-40e3-454f-bd23-f5ec2fe4bac8"
  plan_unique_identifier = "application-runtime-memory"
}

resource "btp_subaccount_entitlement" "entitlement_2" {
  plan_name              = "free"
  service_name           = "content-agent-ui"
  subaccount_id          = "b9b11ec2-40e3-454f-bd23-f5ec2fe4bac8"
  plan_unique_identifier = "content-agent-ui-free"
}

resource "btp_subaccount_entitlement" "entitlement_1" {
  plan_name              = "default"
  service_name           = "auditlog-management"
  subaccount_id          = "b9b11ec2-40e3-454f-bd23-f5ec2fe4bac8"
  plan_unique_identifier = "auditlog-management-default"
}
