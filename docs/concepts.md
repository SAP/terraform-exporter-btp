# How the btptf CLI Works

Here are the steps the btptf CLI takes over for you

1. **Resource Identification**: It creates Terraform `list` blocks and runs `terraform query` against the SAP BTP and official Cloud Foundry providers.
2. **Import Process**: Terraform generates an identity-based import block for every selected query result.
3. **Configuration Generation**: The same query generates the Terraform resource configuration in HashiCorp Configuration Language (HCL).
4. **Configuration Refinements**
It not only creates the Terraform configuration based on the data available on SAP BTP, but also refines the code according to good practice for good readability and sustainable code management. See [How the btptf CLI Refines the Generated Configurations](tfcodeimprovements.md).

!!! info
    No state file is created by the btptf CLI. The reason is that we want to enable best practices and allow the user to add a remote state storage configuration (always customer specific) to the configuration before triggering the state import.
