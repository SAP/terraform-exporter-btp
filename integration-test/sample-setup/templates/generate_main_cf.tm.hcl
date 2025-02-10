generate_hcl "_terramate_generated_main.tf" {
  condition = tm_contains(terramate.stack.tags, "cloudfoundry")
  content {

    resource "cloudfoundry_space" "project_space" {
      name = lower(replace("${tm_upper(terramate.stack.tags[1])}-${var.project_name}", " ", "-"))
      org  = var.cf_org_id
    }

    resource "cloudfoundry_space_role" "space_manager" {
      count    = tm_ternary(terramate.stack.tags[1] == "prod", 0, 1)
      username = var.cf_space_manager
      type     = "space_manager"
      space    = cloudfoundry_space.project_space.id
      origin   = "sap.ids"
    }

    resource "cloudfoundry_space_role" "space_developer" {
      count    = tm_ternary(terramate.stack.tags[1] == "prod", 0, 1)
      username = var.cf_space_developer
      type     = "space_developer"
      space    = cloudfoundry_space.project_space.id
      origin   = "sap.ids"
    }

    resource "cloudfoundry_space" "project_space2" {
      name = lower(replace("${tm_upper(terramate.stack.tags[1])}-test-space", " ", "-"))
      org  = var.cf_org_id
    }
  }
}
