// TERRAMATE: GENERATED AUTOMATICALLY DO NOT EDIT

resource "cloudfoundry_space" "project_space" {
  name = lower(replace("DEV-${var.project_name}", " ", "-"))
  org  = var.cf_org_id
}
resource "cloudfoundry_space_role" "space_manager" {
  count    = 1
  origin   = "sap.ids"
  space    = cloudfoundry_space.project_space.id
  type     = "space_manager"
  username = var.cf_space_manager
}
resource "cloudfoundry_space_role" "space_developer" {
  count    = 1
  origin   = "sap.ids"
  space    = cloudfoundry_space.project_space.id
  type     = "space_developer"
  username = var.cf_space_developer
}
resource "cloudfoundry_space" "project_space2" {
  name = lower(replace("DEV-test-space", " ", "-"))
  org  = var.cf_org_id
}
