terraform {
  required_providers {
    btp = {
      source  = "SAP/btp"
      version = "1.10.0"
    }
  }

  backend "azurerm" {
    resource_group_name  = "rg-terraform-state"
    storage_account_name = "terraformstatestorage"
    container_name       = "tfstate"
    key                  = "terraform.tfstate"
  }
}

provider "btp" {
  globalaccount = var.globalaccount
}
