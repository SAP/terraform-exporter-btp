resource "btp_subaccount_service_instance" "serviceinstance_0" {
  name           = "dev-inttest-exporter-alert-notification"
  serviceplan_id = data.btp_subaccount_service_plan.alert-notification_standard.id
  shared         = false
  subaccount_id  = btp_subaccount.subaccount_0.id
}
