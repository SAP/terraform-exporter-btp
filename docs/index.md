# Terraform Exporter for SAP BTP

The *Terraform Exporter for SAP BTP* (btptf CLI) is a handy tool that makes it easier to bring your existing SAP Business Technology Platform ([BTP](https://www.sap.com/products/technology-platform/what-is-sap-business-technology-platform.html)) resources into [Terraform](https://www.terraform.io/). With it, you can take things like subaccounts and directories in BTP and turn them into Terraform state and configuration files. It's especially useful for teams who are moving to Terraform but still need to manage older infrastructure or SAP BTP accounts that are already set up.


## How does it work

- **Resource Identification**: Terraform Exporter for SAP BTP identifies the SAP BTP resources and maps them to corresponding Terraform resources using the BTP CLI Server APIs.
- **Import Process**: The tool utilizes Terraform's [import](https://developer.hashicorp.com/terraform/cli/import) function to integrate each resource into Terraform's state.
- **Configuration Generation**: After import, it generates the Terraform code (in HashiCorp Configuration Language - HCL) for each resource, enabling further customizations as needed.


