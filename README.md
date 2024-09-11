![Golang](https://img.shields.io/badge/Go-1.23-informational)
[![REUSE status](https://api.reuse.software/badge/github.com/SAP/terraform-exporter-for-sap-btp)](https://api.reuse.software/info/github.com/SAP/terraform-exporter-for-sap-btp)

# Terraform exporter for SAP BTP

## About this project

The *Terraform Exporter for SAP BTP* is a tool that helps export resources from a BTP Global Account.  It can generate Terraform scripts for the resources and import those resources into a Terraform state file.

## Setup

You have two options to setup the CLI:

1. Local build
1. Download the pre-built binary

The following sections describe the details for the two options.


### Pre-built binary

The easiest way to get the binary is to download from the [releases section](https://github.com/SAP/terraform-exporter-btp/releases) of this repository. Select the version you want to use and download the binary that fits your operating system from the `assets` of the release. We recommend using the latest version.


### Local build

If you want to build the binary from scratch, follow these steps:

1. Open this repository inside VS Code Editor
1. We have setup a devcontainer, so reopen the repository in the devcontainer.
1. Open a terminal in VS Code and install the binary by running

   ```bash
    make install
    ```
   This will implicitly trigger a build of the source. If you want to build *without* install, execute `make build`.

1. The system will store the binary as `btptfexporter` (`btptfexporter.exe` in case of Windows) in the default binary path of your Go installation `$GOPATH/bin`.

   > **Note** - You find the value of the GOPATH via `go env GOPATH`

#### Troubleshooting

##### Binary not executable (MacOS or Linux)

In case you get an error that the binary is not executable, naviigate to the location of the binary and execute the following command:

```bash
chomd +x btptfexporter
```

## Usage

After executing the [setup](#setup) of the CLI you must set some required environment variables needed for authentication.

1. Set the environment variable `BTP_GLOBALACCOUNT` which specifies the *subdomain* of your SAP BTP global account.
1. Depending on the authentication flow, set the following environment variables:

   - Basic Authentication: set the environment variable `BTP_USERNAME` and `BTP_PASSWORD`
   - X509 Authentication: set the environment variables `BTP_TLS_CLIENT_CERTIFICATE`, `BTP_TLS_CLIENT_KEY`, `BTP_TLS_IDP_URL`

1. In addition you can set the following optional parameters as environment variables, depending on your requirements:

   - Specify a custom IdP for the authentication via `BTP_IDP`
   - Specify a URL of the BTP CLI server (SAP internal only) via `BTP_CLI_SERVER_URL`

The parameters correspond to the Terraform provider configuration options you find in the [BTP Terraform Provider documentation](https://registry.terraform.io/providers/SAP/btp/latest/docs)

How to set the parameters depends on your setup and is OS-specific:

- On Windows (example):

   ```powershell
   $env:BTP_USERNAME=<MY SAP BTP USERNAME>
   ```

- On MacOS and Linux (example):

   ```bash
   export BTP_USERNAME=<MY SAP BTP USERNAME>
   ```

- In a devcontainer:
   - Create a file `devcontainer.env` in the `.devcontainer` directory
   - Add the environment variables in the file. Here is an example:

      ```txt
      BTP_USERNAME='<MY SAP BTP USERNAME>'
      BTP_PASSWORD='<MY SAP BTP PASSWORD>'
      BTP_GLOBALACCOUNT='<MY SAP BTP GLOBAL ACCOUNT SUBDOMAIN>'
      ```
  - Start the devcontainer variant `Terraform exporter for SAP BTP - Development (with env file)`. The environment variables defined in the .`devcontainer.env` file will be automatically injected.

- Alternative via `.env` file (available on MacOS and Linux only):
   - Create a file `.env` in the root of the project
   - Add the environment variables in the file. Here is an example:

      ```txt
      BTP_USERNAME='<MY SAP BTP USERNAME>'
      BTP_PASSWORD='<MY SAP BTP PASSWORD>'
      BTP_GLOBALACCOUNT='<MY SAP BTP GLOBAL ACCOUNT SUBDOMAIN>'
      ```

    - Execute the following command in a terminal:

       ```bash
       export $(xargs <.env)`
       ```

    > **Note** - There is no predefined fucntionality in PowerShell to achieve the same. A custom script would be needed.

## Commands

The CLI offers several commands for the export of Terraform configurations of SAP BTP. You find a comprehensive overview of the commands and the options in the [documentation](./docs/btptfexporter.md).

## Developer Guide

If you want to contribute to the code of the Terraform Exporter for SAP BTP, please check our [Contribution Guidelines](CONTRIBUTING.md). The technical setup and how to get started are described in the [Developer Guide](DEVELOPER-GUIDE.md)

## Support, Feedback, Contributing

This project is open to feature requests/suggestions, bug reports etc. via [GitHub issues](https://github.com/SAP/terraform-exporter-for-sap-btp/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Security / Disclosure
If you find any bug that may be a security problem, please follow our instructions at [in our security policy](https://github.com/SAP/terraform-exporter-for-sap-btp/security/policy) on how to report it. Please do not create GitHub issues for security-related doubts or problems.

## Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](https://github.com/SAP/.github/blob/main/CODE_OF_CONDUCT.md) at all times.

## Licensing

Copyright 2024 SAP SE or an SAP affiliate company and terraform-exporter-for-sap-btp contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/SAP/terraform-exporter-for-sap-btp).
