# Limitations

## Supported Resources for Import

The btptf CLI can create import blocks and the corresponding configurations only for resources that support the import functionality of Terraform. Not all resources available in the Terraform providers support this feature and can hence not be imported.

You find a list of supported resources for the Terraform Provider for SAP BTP in the corresponding repository on GitHub under the [Overview on importable resources](https://github.com/SAP/terraform-provider-btp/blob/main/guides/IMPORT.md).

## Restrictions for Cloud Foundry

The btptf CLI can create import blocks and the corresponding configurations only for resources that support the import functionality of Terraform. Not all resources available in the Terraform providers support this feature and can hence not be imported. For details please check the [documentation](https://registry.terraform.io/providers/cloudfoundry/cloudfoundry/latest).

In case of resources for Cloud Foundry the btptf CLI focuses on the resources that are available in the Cloud Foundry environment on SAP BTP. It is not intended to be a generic tool for vanilla Cloud Foundry deployments.

You find the details about supported and unsupported Cloud Foundry features on SAP BTP on [help.sap.com](https://help.sap.com/docs/btp/sap-business-technology-platform/cloud-foundry-environment#supported-and-unsupported-cloud-foundry-features).

## Temporal Limitations

There might be some temporal limitations as we rely on the functions and features of the Terraform providers for SAP BTP and Cloud Foundry. These limitations might be removed in an upcoming release of the providers. In this section we list these limitations to avoid confusion and to be transparent what might not yet work with the Terraform Exporter for SAP BTP.

- Export of users for Cloud Foundry: Currently there is a limitations for the resource `cloudfoundry_users` which cannot be exported for a Cloud Foundry environment on SAP BTP. This is also documented in the [issue 298](https://github.com/SAP/terraform-exporter-btp/issues/298). A fix is planned for the Terraform provider for Cloud Foundry. After the availability of the fix we can update the Terraform Exporter for SAP BTP and support this resource.
