generate_hcl "_terramate_generated_provider.tf" {
  condition = tm_contains(terramate.stack.tags, "cloudfoundry")
  content {
    terraform {
      required_providers {
        cloudfoundry = {
          source  = "cloudfoundry/cloudfoundry"
          version = tm_ternary(tm_contains(terramate.stack.tags, "dev"), global.terraform.providers.cloudfoundry.version_dev, global.terraform.providers.cloudfoundry.version)
        }
      }
    }
    provider "cloudfoundry" {
      api_url = var.cf_api_url
    }
  }
}
