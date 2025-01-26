// TERRAMATE: GENERATED AUTOMATICALLY DO NOT EDIT

output "cf_api_url" {
  value     = jsondecode(btp_subaccount_environment_instance.cloudfoundry.labels)["API Endpoint"]
  sensitive = false
}
output "cf_org_id" {
  value     = jsondecode(btp_subaccount_environment_instance.cloudfoundry.labels)["Org ID"]
  sensitive = false
}
