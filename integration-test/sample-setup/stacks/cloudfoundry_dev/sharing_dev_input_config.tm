input "cf_api_url" {
  backend       = "default"
  from_stack_id = "cc4ad8b2-3d7e-4719-82e5-de9d2a299611"
  value         = outputs.cf_api_url.value
  mock          = "https://api.cf.us21.hana.ondemand.com"
}

input "cf_org_id" {
  backend       = "default"
  from_stack_id = "cc4ad8b2-3d7e-4719-82e5-de9d2a299611"
  value         = outputs.cf_org_id.value
  mock          = "917f57a1-8fee-43b3-b3a8-4bb4ce8259ab"
}
