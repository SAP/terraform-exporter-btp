# Limitations

## Supported SAP BTP Resources

The btptf CLI can create import blocks and the corresponding configurations only for resources that support the import functionality of Terraform. Not all resources available in the Terraform providers support this feature and can hence not be imported.

You find a list of supported resources for the Terraform Provider for SAP BTP in the corresponding repository on GitHub under the [Overview on importable resources](https://github.com/SAP/terraform-provider-btp/blob/main/guides/IMPORT.md).

## Supported Cloud Foundry Resources

The btptf CLI can create import blocks and the corresponding configurations only for resources that support the import functionality of Terraform. Not all resources available in the Terraform providers support this feature and can hence not be imported. For details please check the [documentation](https://registry.terraform.io/providers/cloudfoundry/cloudfoundry/latest).

The btptf CLI focuses on the resources that are available in the Cloud Foundry environment on SAP BTP. It is not intended to be a generic tool for vanilla Cloud Foundry deployments.

You find the details about supported and unsupported Cloud Foundry features on SAP BTP on [help.sap.com](https://help.sap.com/docs/btp/sap-business-technology-platform/cloud-foundry-environment#supported-and-unsupported-cloud-foundry-features).

