# Remote Backend Configuration

The Terraform Exporter for SAP BTP generates the configuration with a [local state backend](https://developer.hashicorp.com/terraform/language/backend/local). This is not recommended for a productive setup. You should use a [remote backend](https://developer.hashicorp.com/terraform/language/backend) for storing the state

You have two options to introduce such a remote backend to your configuration:

1. After the configuration code is generated you add the corresponding block in the `provider.tf` file
2. You inject the backend configuration during the code generation using the CLI flags `--backend-path`to specify the path to a sample backend configuration or `--backend-type` and `--backend-config` to explicitly specify the parameters

!!! info
    The Terraform Exporter for SAP BTP executes the `terraform init` with the option `-backend=false` independent if a backend was configured or not to make sure that the basic initialization has taken place.

## Examples

The following sections showcase examples of injecting the backend configuration when the configuration is generated. We assume that we want to use a Azure Blob Storage as state backend.

### Remote State Store via a Sample File

We specify the remote state store via a sample file `backend.tf` that has the following layout:

```terraform
terraform {
  backend "azurerm" {
    resource_group_name  = "rg-terraform-state"
    storage_account_name = "terraformstatestorage"
    container_name       = "tfstate"
    key                  = "terraform.tfstate"
  }
}
```

Be aware that the `terraform` block is mandatory when using the sample file.

The file is stored in the same directory that the CLI will be called from. To instruct the CLI to use this file we must use the  `--backend-path` or short `-b` flag and set the path to the file. The CLI command for the export of a Terraform configuration injecting this sample backend configuration looks like this:

```bash
btptf export -s 23fe9a1b-923d-4ab0-ae24-86ff3384cf93 -b backend.tf
```

### Remote State Store via a CLI Parameters

We want to specify the different parameters of the state backend as flags of the CLI execution. The flags that we will use are:

- `--backend-type` or short `-t`to specify the backend type
- `--backend-config` or short `-e` to specify the parameters.This flag can be used multiple times.

The CLI command for the export of a Terraform configuration injecting this backend configuration via the flags looks like this:

```bash
btptf export -r='subaccount' -s 23fe9a1b-923d-4ab0-ae24-86ff3384cf93 -t azurerm \
-e 'resource_group_name=rg-terraform-state' \
-e 'storage_account_name=terraformstatestorage' \
-e 'container_name=tfstate' \
-e 'key=terraform.tfstate'
```
