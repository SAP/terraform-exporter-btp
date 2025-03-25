# Terraform Exporter for SAP BTP

## Introduction
Welcome to the documentation for the Terraform exporter for SAP BTP (**btptf CLI**). It's a tool that helps you integrate existing SAP BTP resources into Terraform, facilitating a seamless adoption of Infrastructure-as-Code practices. 

This guide will walk you through the concepts, features, and benefits of using the btptf CLI, and it provides comprehensive how-tos.

The following SAP BTP account levels are available for export either as complete entities or by choosing particular resources individually: 

- Directories
- Subaccounts
- Cloud Foundry orgs


!!! info
    The Terraform Exporter for SAP BTP is fully compatible with [OpenTofu](https://opentofu.org/). All steps outlined in this guide can be executed using either the Terraform CLI or the OpenTofu CLI. For simplicity, this documentation will only reference [Terraform](https://www.terraform.io/).


## Benefits

The btptf CLI offers several advantages that streamline the adoption of Infrastructure-as-Code for SAP BTP with Terraform:


- **Enhanced User Experience**
 There is a native import functionality of Terraform that is very generic and implies a lot of manual work when importing resources from SAP BTP into Terraform. The btptf CLI streamlines this import process specifically for SAP BTP by combining multiple individual steps into one cohesive workflow. 
For an example of how the import for SAP BTP resources works without the btptf CLI, just with the Terraform import commands, you can explore [this sample repository](https://github.com/SAP-samples/btp-terraform-samples/tree/main/released/import).

- **Automated Resource Import**
The btptf CLI automatically converts BTP resources into Terraform code, minimizing manual adjustments required by the generic Terraform CLI import.

- **Bulk Resource Generation**
Users can generate Terraform code for multiple resources with one command, specifying which resources to include for tailored exports. Per export command, users need to specify an SAP BTP account level as execution level: directoriy, subaccount, or Cloud Foundry org.

- **Quality Code Practices**
The btptf CLI ensures clean, maintainable Terraform code generation, adhering to good practices for readability and efficiency.

- **Flexible JSON Workflow or Direct Export**
The btptf CLI offers a JSON workflow where a JSON-based resource inventory is created that users can edit before export, allowing pre-export customization. Users can also directly export resources.









