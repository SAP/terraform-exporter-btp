resource "btp_subaccount_subscription" "subscription_5" {
  app_name      = "feature-flags-dashboard"
  parameters    = jsonencode({})
  plan_name     = "dashboard"
  subaccount_id = "1234567890"
  timeouts      = null
}
