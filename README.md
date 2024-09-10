![Golang](https://img.shields.io/badge/Go-1.23-informational)
[![REUSE status](https://api.reuse.software/badge/github.com/SAP/terraform-exporter-for-sap-btp)](https://api.reuse.software/info/github.com/SAP/terraform-exporter-for-sap-btp)

# Terraform exporter for SAP BTP

## About this project

The *Terraform Exporter for SAP BTP* is a tool that helps export resources from a BTP Global Account.  It can generate Terraform scripts for the resources and import those resources into a Terraform state file.

## Requirements and Setup

### Local build

If you want to build the binary from scratch, folow these steps:

1. Open this repo inside VS Code Editor
1. We have setup a devcontainer, so open the repo using devcontainer.
1. Build the binary: From the terminal in vscode run `make build` & `make install`
1. A file (binary) `btptfexporter` will be found in the default binary path of your Go installation.
1. Make it executable: `chomd +x btptfexporter`.

### Prebuilt binary

Please go to the releases section of this repository and download the binary for your system.

## Usage

1. [Download](https://github.tools.sap/BTP-Terraform/btptfexporter/releases/tag/v0.0.3-poc) or build the binary to a local path/folder.
1. Create the following required environment varaibles: `BTP_USERNAME`, `BTP_PASSWORD`, `BTP_GLOBALACCOUNT`
Optionally, you can set the following parameters: `BTP_CLIENT_SERVER_URL`, `BTP_IDP`, `BTP_TLS_CLIENT_CERTIFICATE`, `BTP_TLS_CLIENT_KEY`, `BTP_TLS_IDP_URL`. Please refer the [BTP Terraform Provider documentation](https://registry.terraform.io/providers/SAP/btp/latest/docs) to learn more about these parameters.

    > **Note** - You can also define a `.env` file with the content of the parameters and export this file via `export $(xargs <.env)`.

1. use the `--help` flag to know more.

## Commands

### resource : Export specific btp resources from a subaccount

Use this command to create terraform configuration for all the resources of a subaccount or specific resource using the subcommands

```bash
btptfexporter resource [command]

Example:

btptfexporter resource all --subaccount <subaccount-id>

Available Commands:

  all                   export all resources of a subaccount
  entitlements          export entitlements of a subaccount
  environment-instances export environment instance of a subaccount
  from-file             export resources from a json file.
  subaccount            export subaccount
  subscriptions         export subscriptions of a subaccount
  trust-configurations  export trust configurations of a subaccount

  ```

### generate-resources-list  : Store the list of resources from btp subaccount into a json file

Use this command to get the list of resources from a subaccont and store it in a json file.

```bash
btptfexporter generate-resources-list [flags]

Example:

btptfexporter generate-resources-list --resources=entitlements,subscriptions --subaccount=<subacount_id>
```

Valid resources are:
- subaccount
- entitlements
- subscriptions
- environment-instances
- trust-configurations


### Generate markdown documentation

We can generate the markdown documentation via the make file:

```bash
make docs
```

## Support, Feedback, Contributing

This project is open to feature requests/suggestions, bug reports etc. via [GitHub issues](https://github.com/SAP/terraform-exporter-for-sap-btp/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Security / Disclosure
If you find any bug that may be a security problem, please follow our instructions at [in our security policy](https://github.com/SAP/terraform-exporter-for-sap-btp/security/policy) on how to report it. Please do not create GitHub issues for security-related doubts or problems.

## Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](https://github.com/SAP/.github/blob/main/CODE_OF_CONDUCT.md) at all times.

## Licensing

Copyright 2024 SAP SE or an SAP affiliate company and terraform-exporter-for-sap-btp contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/SAP/terraform-exporter-for-sap-btp).
