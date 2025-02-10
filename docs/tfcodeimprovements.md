# Terraform Configuration Improvements

!!! info
    As of now this is an experimental feature and needs to be activated by setting the environment varibale `BTPTF_EXPERIMENTAL` to a non-empty value.

The Terraform Exporter for SAP BTP creates the Terraform configuration based on the data available on SAP BTP. This means that the resulting configuration needs manual adjustements before beeing in a decent state.

We are continously improving this situation by refining the configuration after it gets generated. The following sections give an overview of the improvements.

## General Configuration Improvements

### Remove empty values

In general we clean up empty values from the configuration namely if attributes are:

- `null`
- emtpy JSON strings

### Extract Provider configuration

We extract the value of the provider configuration (`provider.tf` file) as a variable namely:

- The subdomain of the SAP BTP global acccount for the Terraform provider for SAP BTP
- The API URL of Cloud Foundry for the Terraform provider for Cloud Foundry

## Improvements on Subaccount Level

The Terraform configuration on subaccount level gets improved via the following measures:

- The `region` is extracted as a variable for the `btp_subaccount` resources.
- All resources that reference a `subaccount_id` get transformed to reference the `btp_subaccount` resource if available.
- The attribute `parent_id` of the resource `btp_subaccount` gets extracted as a variable.
- If a dependency between the resource `btp_subaccount_entitlement` and the resource `btp_subaccount_subscription` exists, a corresponding `depends_on` block will be added to the resource `btp_subaccount_subscription`.

## Improvements on Directory Level

- The attribute `parent_id` of the resource `btp_directory` gets extracted as a variable.
- All resources that reference a `directory_id` get transformed to reference the `btp_directory` resource if available.

## Improvements on Cloud Foundry Organizational Level

No improvements available yet
