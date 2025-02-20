resource "btp_subaccount_subscription" "subscription_5" {
  app_name      = "feature-flags-dashboard"
  plan_name     = "dashboard"
  subaccount_id = "1234567890"
  depends_on = [ btp_subaccount_entitlement.entitlement_2 ]
}
