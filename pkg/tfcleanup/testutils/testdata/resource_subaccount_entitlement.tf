resource "btp_subaccount_entitlement" "entitlement_2" {
  amount        = 1
  plan_name     = "dashboard"
  service_name  = "feature-flags-dashboard"
  subaccount_id = "1234567890"
}
