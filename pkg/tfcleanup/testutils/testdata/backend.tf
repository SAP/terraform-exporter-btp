terraform {
  backend "azurerm" {
    resource_group_name  = "rg-terraform-state"
    storage_account_name = "terraformstatestorage"
    container_name       = "tfstate"
    key                  = "terraform.tfstate"
  }
}
