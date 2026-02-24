
###
# Resource: BTP_SUBACCOUNT_ENTITLEMENT
###
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform
resource "btp_subaccount_entitlement" "entitlement_3" {
  plan_name     = "api_access"
  service_name  = "connectivity"
  subaccount_id = "b9b11ec2-40e3-454f-bd23-f5ec2fe4bac8"
}

# __generated__ by Terraform
resource "btp_subaccount_entitlement" "entitlement_4" {
  plan_name     = "standard"
  service_name  = "cloudfoundry"
  subaccount_id = "b9b11ec2-40e3-454f-bd23-f5ec2fe4bac8"
}

# __generated__ by Terraform
resource "btp_subaccount_entitlement" "entitlement_0" {
  amount        = 2
  plan_name     = "MEMORY"
  service_name  = "APPLICATION_RUNTIME"
  subaccount_id = "b9b11ec2-40e3-454f-bd23-f5ec2fe4bac8"
}

# __generated__ by Terraform
resource "btp_subaccount_entitlement" "entitlement_2" {
  plan_name     = "free"
  service_name  = "content-agent-ui"
  subaccount_id = "b9b11ec2-40e3-454f-bd23-f5ec2fe4bac8"
}

# __generated__ by Terraform
resource "btp_subaccount_entitlement" "entitlement_1" {
  plan_name     = "default"
  service_name  = "auditlog-management"
  subaccount_id = "b9b11ec2-40e3-454f-bd23-f5ec2fe4bac8"
}
