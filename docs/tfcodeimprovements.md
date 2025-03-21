# Terraform Configuration Improvements

The Terraform Exporter for SAP BTP creates the Terraform configuration based on the data available on SAP BTP. This means that the resulting configuration needs manual adjustments before being in a decent state.

We are constantly improving this situation by refining the configuration after it gets generated. The following sections give an overview of the improvements.

## General Configuration Improvements

### Remove empty values

In general we clean up empty values from the configuration namely if attributes are:

- `null`
- empty JSON strings

### Extract Provider configuration

We extract the value of the provider configuration (`provider.tf` file) as a variable namely:

- The subdomain of the SAP BTP global account for the Terraform provider for SAP BTP
- The API URL of Cloud Foundry for the Terraform provider for Cloud Foundry

### Remove automatically provided resources

We remove some resources that get automatically created by SAP BTP when e.g., a subaccount or directory get created.

The following resources get removed if they can be identified as automatically created:

- Default trust configuration to the SAP IdP
- Default role collections that get created during subaccount or directory creation or when a service instance or app subscription gets created
- Default roles that get created during subaccount or directory creation or when a service instance or app subscription gets created

## Improvements on Subaccount Level

The Terraform configuration on subaccount level gets improved via the following measures:

- The `region` is extracted as a variable for the resource `btp_subaccount`.
- All resources that reference a `subaccount_id` get transformed to reference the resource `btp_subaccount` if available.
- If the attribute `parent_id` is the global account, it gets removed from the resource `btp_subaccount`. If not it gets extracted as a variable.
- If a dependency between the resource `btp_subaccount_entitlement` and the resource `btp_subaccount_subscription` exists, a corresponding `depends_on` block will be added to the resource `btp_subaccount_subscription`.
- If a dependency between the resource `btp_subaccount_entitlement` and the resource `btp_subaccount_service_instance` exists, a data source `btp_subaccount_service_plan` is added to fetch the identifier of the service plan. In addition a `depends_on` block is added to the data source to make it explicitly dependent to the corresponding entitlement. Finally the technical identifier of the service plan is replaced by the reference to the data source
- If a resource `btp_subaccount_trust_configuration` gets exported that is the SAP default trust, the corresponding resource is removed from the resource configuration as well as the import block from the corresponding import file.
- If a resource `btp_subaccount_role_collection` as well as `btp_subaccount_role` a check if there are dependencies is executed. If there are any a corresponding `depends_on` block is added to the resource `btp_subaccount_role_collection`.

## Improvements on Directory Level

- If the attribute `parent_id` is the global account, it gets removed from the resource `btp_directory`. If not it gets extracted as a variable.
- All resources that reference a `directory_id` get transformed to reference the `btp_directory` resource if available.
- If a resource `btp_directory_role_collection` as well as `btp_directory_role` a check if there are dependencies is executed. If there are any a corresponding `depends_on` block is added to the resource `btp_directory_role_collection`.


## Improvements on Cloud Foundry Organizational Level

- The ID of the organization is extracted as a variable
