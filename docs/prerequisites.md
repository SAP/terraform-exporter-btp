# Prerequisites

- You must have one of the following CLIs installed:
     - [Terraform CLI](https://developer.hashicorp.com/terraform/install).
     - [OpenTofu CLI](https://opentofu.org/docs/intro/install/).

    If you download the CLI binaries make sure that they are in your `PATH` environment variable.
    If you have both CLIs installed, you can enforce the usage of one via the environment variable `BTPTF_IAC_TOOL` by setting its value to `terraform` or `tofu` respectively. Otherwise the btptf CLI will check for Terraform first and use it i available.

- To export directories or subaccounts, you need to set the following environment variables to authenticate against SAP BTP:

| Environment Variable Name  | Description |
| --- | --- |
| BTP_GLOBALACCOUNT | The subdomain of the global account from which you want to import resources. |
| BTP_USERNAME | Your user name, usually an e-mail address. |
| BTP_PASSWORD | Your password. Note that two-factor authentication is not supported.  |
| BTP_TLS_CLIENT_CERTIFICATE | PEM encoded certificate (only required for x509 authentication) |
| BTP_TLS_CLIENT_KEY | PEM encoded private key (only required for x509 authentication) |
| BTP_TLS_IDP_URL | The URL of the identity provider to be used for authentication (only required for x509 authentication) |
| BTP_IDP | The identity provider to be used for authentication (only required for custom IDP) |
| BTP_CLI_SERVER_URL | The URL of the BTP CLI server (Relevant for SAP internal use-cases only)  |

!!! Warning

    Do not set the `BTP_ENABLE_SSO` parameter when using the btptf CLI. Processing will abort as this parameter is not supported


The parameters are the ones required by the [BTP Terraform Provider](https://registry.terraform.io/providers/SAP/btp/latest/docs).

- To export directories or subaccounts, you need global account administrator permissions.

- To export Cloud Foundry orgs, you need to authenticate against Cloud Foundry by setting the following environment variables:

| Environment Variable Name  | Description |
| --- | --- |
| CF_API_URL | Specific URL representing the entry point for communication between the client and a Cloud Foundry instance. |
| CF_USER | A unique identifier associated with an individual or entity for authentication & authorization purposes. |
| CF_PASSWORD | A confidential alphanumeric code associated with a user account on the Cloud Foundry platform, requires user to authenticate.  |
| CF_ORIGIN | Indicates the identity provider to be used for login |
| CF_CLIENT_ID | Unique identifier for a client application used in authentication and authorization processes. |
| CF_CLIENT_SECRET | A confidential string used by a client application for secure authentication and authorization, requires cf_client_id to authenticate |
| CF_ACCESS_TOKEN | OAuth token to authenticate with Cloud Foundry |
| CF_REFRESH_TOKEN | Token to refresh the access token, requires access_token |

These environment variables are the ones required by the [Terraform provider for Cloud Foundry](https://registry.terraform.io/providers/cloudfoundry/cloudfoundry/latest/docs).

- To export Cloud Foundry orgs, you need the Org Manager role.

## How to set the parameters
Depending on your operating systems, you set the environment variables as follows:

=== "Windows"

    ``` powershell
    $env:BTP_USERNAME=<MY SAP BTP USERNAME>
    ```

=== "Linux/Mac"

    ``` bash
    export BTP_USERNAME=<MY SAP BTP USERNAME>
    ```
