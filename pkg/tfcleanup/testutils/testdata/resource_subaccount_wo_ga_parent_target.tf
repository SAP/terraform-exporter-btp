resource "btp_subaccount" "subaccount_0" {
  beta_enabled = true
  description  = null
  labels = {
    costcenter = ["12345"]
    stage      = ["DEV"]
  }
  name      = "DEV Inttest Exporter (DO NOT DELETE)"
  parent_id = var.parent_id
  region    = "us10"
  subdomain = "dev-inttest-exporter-8b2967c9-c3f3-68a4-4604-16702098e1d1"
  usage     = "NOT_USED_FOR_PRODUCTION"
}
