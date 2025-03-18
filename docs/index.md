# Terraform Exporter for SAP BTP

The *Terraform Exporter for SAP BTP* (btptf CLI) exports existing SAP BTP resources as Terraform code, so you can start adopting Infrastructure-as-Code with Terraform.

The following SAP BTP account levels can be exported:

- Directories
- Subaccounts
- Cloud Foundry environment instances

!!! info
    The Terraform Exporter for SAP BTP is fully compatible with [OpenTofu](https://opentofu.org/). All steps outlined in this guide can be executed using either the Terraform CLI or the OpenTofu CLI. For simplicity, this documentation will only reference [Terraform](https://www.terraform.io/).


## How does it work

- **Resource Identification**: Terraform Exporter for SAP BTP identifies the SAP BTP resources and maps them to corresponding Terraform resources using the BTP CLI Server APIs.
- **Import Process**: The tool utilizes Terraform's [import](https://developer.hashicorp.com/terraform/cli/import) function to integrate each resource into Terraform's state.
- **Configuration Generation**: After import, it generates the Terraform code (in HashiCorp Configuration Language - HCL) for each resource, enabling further customizations as needed.
