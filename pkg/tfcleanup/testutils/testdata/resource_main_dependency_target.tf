resource "btp_subaccount_subscription" "subscription_5" {
  app_name      = "feature-flags-dashboard"
  parameters    = jsonencode({})
  plan_name     = "dashboard"
  subaccount_id = btp_subaccount_dummy.subaccount_0.id
  timeouts      = null
}
