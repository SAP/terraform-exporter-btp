# How the btptf CLI Refines the Generated Configurations

The btptf CLI not only creates the Terraform configuration based on the data available on SAP BTP, but it also cleans up the the resulting configurations by refining the code and adding a `variables.tf` file.

The following section outlines the code and, more generally,  configuration refinements carried out by the btptf CLI.

## General Refinements

### Remove empty values

Empty values are removed from the code. In particular: attributes:

- `null`
- empty JSON strings

### Extract provider configuration

The value of the provider configuration (`provider.tf` file) is extracted and stored as a variable:

- The subdomain of the SAP BTP global account for the Terraform provider for SAP BTP
- The API URL of Cloud Foundry for the Terraform provider for Cloud Foundry

### Remove automatically provided resources

Some superfluous resources that get automatically created by SAP BTP at, e.g. subaccount or directory creation, are removed:

- Default trust configuration to the SAP IdP
- Default role collections that get created during subaccount or directory creation or when a service instance or app subscription gets created
- Default roles that get created during subaccount or directory creation or when a service instance or app subscription gets created
- Default entitlements that get created during the initial subaccount creation

## Refinements on Subaccount Level

The Terraform configuration on subaccount level gets improved via the following measures:

- The `region` is extracted as a variable for the resource `btp_subaccount`.
- All resources that reference a `subaccount_id` get transformed to reference the resource `btp_subaccount` if available.
- If the attribute `parent_id` is the global account, it gets removed from the resource `btp_subaccount`, as the global account is the default. If not it gets extracted as a variable.
- If a dependency between the resource `btp_subaccount_entitlement` and the resource `btp_subaccount_subscription` exists, a corresponding `depends_on` block will be added to the resource `btp_subaccount_subscription`.
- If a dependency between the resource `btp_subaccount_entitlement` and the resource `btp_subaccount_service_instance` exists, a data source `btp_subaccount_service_plan` is added to fetch the identifier of the service plan. In addition a `depends_on` block is added to the data source to make it explicitly dependent to the corresponding entitlement. Finally the technical identifier of the service plan is replaced by the reference to the data source`.
- If a resource `btp_subaccount_trust_configuration` gets exported that is the SAP default trust, the corresponding resource is removed from the resource configuration as well as the import block from the corresponding import file.
- If a resource `btp_subaccount_role_collection` as well as `btp_subaccount_role` a check if there are dependencies is executed. If there are any a corresponding `depends_on` block is added to the resource `btp_subaccount_role_collection`.

## Refinements on Directory Level

- If the attribute `parent_id` is the global account, it gets removed from the resource `btp_directory`, as the global account is the default. If not it gets extracted as a variable.
- All resources that reference a `directory_id` get transformed to reference the `btp_directory` resource if available.
- If a resource `btp_directory_role_collection` as well as `btp_directory_role` a check if there are dependencies is executed. If there are any a corresponding `depends_on` block is added to the resource `btp_directory_role_collection`.

## Refinements on Cloud Foundry Org Level

- The ID of the org is extracted as a variable

## Refinements in Beta

We continuously improve the generated Terraform configurations. The following refinements are currently in beta and must be explicitly activated via environment variables.

### Entitlement Module Generation on Subaccount Level

By default, the generated code for entitlements contains one resource `btp_subaccount_entitlement` per entitlement. Depending on the number of entitlements the Terraform configuration can become quite extensive. Therefore, some users rely on the **community** module [sap-btp-entitlements](https://registry.terraform.io/modules/aydin-ozcan/sap-btp-entitlements/btp/latest) to make the setup more comprehensive.

Based on a feature request by these users, we introduced the beta feature to adjust the generated code namely to use the community module instead of the individual resources. The following adjustments are made to the generated code:

- All resources `btp_subaccount_entitlement` are replaced by the module `sap-btp-entitlements`
- A variable is generated for the module `sap-btp-entitlements` to hold the entitlements
- Dependencies to the entitlements from subscriptions or service instances are pointing to the module `sap-btp-entitlements` instead of the individual resources
- The import blocks namely the target address of the import are adjusted to point to the module `sap-btp-entitlements`

To activate this feature you must set the environment variable `BTPTF_ADD_ENTITLEMENTMODULE` to any value as follows:

=== "Windows"

    ``` powershell
    $env:BBTPTF_ADD_ENTITLEMENTMODULE='true'
    ```

=== "Linux/Mac"

    ``` bash
    export BTPTF_ADD_ENTITLEMENTMODULE=true
    ```
