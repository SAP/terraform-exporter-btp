// TERRAMATE: GENERATED AUTOMATICALLY DO NOT EDIT

terraform {
  required_providers {
    cloudfoundry = {
      source  = "cloudfoundry/cloudfoundry"
      version = "~> 1.4.0"
    }
  }
}
provider "cloudfoundry" {
  api_url = var.cf_api_url
}
