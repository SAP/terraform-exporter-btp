resource "btp_subaccount_subscription" "subscription_5" {
  app_name      = "feature-flags-dashboard"
  plan_name     = var.plan_name
  subaccount_id = "123456789"
}
