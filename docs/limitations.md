# Limitations

## Supported SAP BTP Resources

The btptf CLI can export resources for which the provider implements Terraform list-resource discovery and configuration generation. Not every provider resource has a corresponding list resource.

You find a list of supported resources for the Terraform Provider for SAP BTP in the corresponding repository on GitHub under the [Overview on importable resources](https://github.com/SAP/terraform-provider-btp/blob/main/guides/IMPORT.md).

## Supported Cloud Foundry Resources

The btptf CLI can export resources for which the official `cloudfoundry/cloudfoundry` provider implements Terraform list-resource discovery and configuration generation. For details, check the [official provider documentation](https://registry.terraform.io/providers/cloudfoundry/cloudfoundry/latest).

The btptf CLI focuses on the resources that are available in the Cloud Foundry environment on SAP BTP. It is not intended to be a generic tool for vanilla Cloud Foundry deployments.

You find the details about supported and unsupported Cloud Foundry features on SAP BTP on [help.sap.com](https://help.sap.com/docs/btp/sap-business-technology-platform/cloud-foundry-environment#supported-and-unsupported-cloud-foundry-features).

## Subaccounts created with `skip_auto_entitlement`

The Terraform provider for SAP BTP is not able to fetch the the information if the subaccount was created with the parameter `skip_auto_entitlement` set to `true`. This flag must be set manually after generating the configuration of the subaccount.

## Import of Destinations

The resource `btp_subaccount_destination_generic` is not currently part of the exporter's supported V2 resource catalog and must be imported manually:

1. Create the [import block](https://developer.hashicorp.com/terraform/language/block/import) for the [destination resource](https://registry.terraform.io/providers/SAP/btp/latest/docs/resources/subaccount_destination_generic#import) in the Terraform configuration.
2. Create the resource configuration for the destination resource in the Terraform configuration. Please make sure to use the same name for the resource as used in the import block. The parameters needed for the configuration can be fetched from the SAP BTP Cockpit. Go to the destination details and click on the "Export" button. Choose as format `JSON` and use the content of the export for the configuration.
3 Execute the import via the `terraform apply` command.
