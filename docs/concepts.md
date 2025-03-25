# How the btptf CLI Works 

Here are the steps the btptf CLI takes over for you 

1. **Resource Identification**: It first identifies your SAP BTP resources, mapping them accurately to corresponding Terraform resources with the help of the BTP CLI Server APIs.
2. **Import Process**: It utilizes Terraform's [import functionality](https://developer.hashicorp.com/terraform/cli/import) to integrate each identified resource into Terraform's state.
3. **Configuration Generation**: After import, it generates the Terraform code (in HashiCorp Configuration Language - HCL) for each resource, enabling further customizations as needed.
4. **Configuration Refinements**
It not only creates the Terraform configuration based on the data available on SAP BTP, but also refines the code according to good practice for good readability and sustainable code management. See [How the btptf CLI Refines the Generated Configurations](tfcodeimprovements.md).




STUFF TO BE MOVED: 

No state file is created by the btptf CLI. The reason is that we want to enable best practices and allow the user to add a remote state storage configuration (always customer specific) to the configuration before triggering the state import.


The configurations delivered by the btptf CLI are:

  -	Provider configuration (excluding credentials)
  -	[Import](https://developer.hashicorp.com/terraform/language/import) blocks for the resources
  -	Resource configuration retrieved from the platform




  