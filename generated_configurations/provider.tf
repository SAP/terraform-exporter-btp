terraform {
required_providers {
btp = {
source  = "SAP/btp"
version = "1.6.0"
}
}
}

provider "btp" {
globalaccount = "terraformintprod"
}