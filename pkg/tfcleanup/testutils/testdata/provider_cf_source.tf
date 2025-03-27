terraform {
  required_providers {
    cloudfoundry = {
      source  = "cloudfoundry/cloudfoundry"
      version = "1.4.0"
    }
  }
}

provider "cloudfoundry" {
  api_url = "https://api.cf.sap.hana.ondemand.com"
}
