# Limitations

## Supported SAP BTP Resources

The btptf CLI can create import blocks and the corresponding configurations only for resources that support the import functionality of Terraform. Not all resources available in the Terraform providers support this feature and can hence not be imported.

You find a list of supported resources for the Terraform Provider for SAP BTP in the corresponding repository on GitHub under the [Overview on importable resources](https://github.com/SAP/terraform-provider-btp/blob/main/guides/IMPORT.md).

## Supported Cloud Foundry Resources

The btptf CLI can create import blocks and the corresponding configurations only for resources that support the import functionality of Terraform. Not all resources available in the Terraform providers support this feature and can hence not be imported. For details please check the [documentation](https://registry.terraform.io/providers/cloudfoundry/cloudfoundry/latest).

The btptf CLI focuses on the resources that are available in the Cloud Foundry environment on SAP BTP. It is not intended to be a generic tool for vanilla Cloud Foundry deployments.

You find the details about supported and unsupported Cloud Foundry features on SAP BTP on [help.sap.com](https://help.sap.com/docs/btp/sap-business-technology-platform/cloud-foundry-environment#supported-and-unsupported-cloud-foundry-features).

## Subaccounts created with `skip_auto_entitlement`

The Terraform provider for SAP BTP is not able to fetch the the information if the subaccount was created with the parameter `skip_auto_entitlement` set to `true`. This flag must be set manually after generating the configuration of the subaccount.

## Import of Destinations

In general the resource `btp_subaccount_destination_generic` (available since release 1.19.0 of the Terraform provider for SAP BTP) supports the import operation. Due to the design of the resource, we cannot use the mechanism of resource configuration generation (see [HashiCorp Configuration language - Generating configuration](https://developer.hashicorp.com/terraform/language/import/generating-configuration)) that is used as central part of this tool. Hence, the import of destinations via this tool cannot be supported and must be done manually. The flow of the manual import is:

1. Create the [import block](https://developer.hashicorp.com/terraform/language/block/import) for the [destination resource](https://registry.terraform.io/providers/SAP/btp/latest/docs/resources/subaccount_destination_generic#import) in the Terraform configuration.
2. Create the resource configuration for the destination resource in the Terraform configuration. Please make sure to use the same name for the resource as used in the import block. The parameters needed for the configuration can be fetched from the SAP BTP Cockpit. Go to the destination details and click on the "Export" button. Choose as format `JSON` and use the content of the export for the configuration.
3 Execute the import via the `terraform apply` command.
