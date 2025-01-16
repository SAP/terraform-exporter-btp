# Prerequisite

## Terraform CLI

The btptf CLI requires a installation of the Terraform CLI. The Terraform CLI will be called by the btptf CLI.

You find the necessary information  in the [official Terraform documentation](https://developer.hashicorp.com/terraform/install#darwin).

## Setting of Environment Variables

!!! info
    A complete overview of available parameters is described in the next section.

After executing the setup of the btptf CLI, you must set some required environment variables needed for authentication. This section describes the minimum set of environment variables that you must set to execute the btptf CLI successfully:

1. Set the environment variable `BTP_GLOBALACCOUNT` which specifies the *subdomain* of your SAP BTP global account.

2. Depending on the authentication flow, set the following environment variables:

    - Basic Authentication: set the environment variable `BTP_USERNAME` and `BTP_PASSWORD`
    - X509 Authentication: set the environment variables `BTP_TLS_CLIENT_CERTIFICATE`, `BTP_TLS_CLIENT_KEY`, `BTP_TLS_IDP_URL`

3. In addition you can set the following optional parameters as environment variables, depending on your requirements:

    - Specify a custom IdP for the authentication via `BTP_IDP`
    - Specify a URL of the BTP CLI server (SAP internal only) via `BTP_CLI_SERVER_URL`
    - Specify the login using SSO via `BTP_ENABLE_SSO` (true/false)

The parameters correspond to the Terraform provider configuration options that you find in the [BTP Terraform Provider documentation](https://registry.terraform.io/providers/SAP/btp/latest/docs)

How to set the parameters depends on your setup and is OS-specific:

=== "Windows"

    ``` powershell
    $env:BTP_USERNAME=<MY SAP BTP USERNAME>
    ```

=== "Linux/Mac"

    ``` bash
    export BTP_USERNAME=<MY SAP BTP USERNAME>
    ```

- In a devcontainer:
    - Create a file `devcontainer.env` in the `.devcontainer` directory
    - Add the environment variables in the file. Here is an example:

      ```txt
      BTP_USERNAME='<MY SAP BTP USERNAME>'
      BTP_PASSWORD='<MY SAP BTP PASSWORD>'
      BTP_GLOBALACCOUNT='<MY SAP BTP GLOBAL ACCOUNT SUBDOMAIN>' #optional
      ```
    - Start the devcontainer option `Terraform exporter for SAP BTP - Development (with env file)`. The environment variables defined in the `devcontainer.env` file will be automatically injected.

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
       export $(xargs <.env)
       ```

!!! info
    There is no predefined functionality in PowerShell to achieve the same. A custom script is needed.

## Overview of Environment Variables

The environment variables supported by the btptf CLI are required to configure the Terraform providers.

### Terraform provider for SAP BTP

For the scenarios where you want to import resources defined in the [Terraform provider for SAP BTP](https://registry.terraform.io/providers/SAP/btp/latest/docs) the following parameters are available:

| Environment Variable Name  | Description |
| --- | --- |
| BTP_GLOBALACCOUNT | The subdomain of the global account from which you want to import resources. |
| BTP_USERNAME | Your user name, usually an e-mail address. |
| BTP_PASSWORD | Your password. Note that two-factor authentication is not supported.  |
| BTP_IDP | The identity provider to be used for authentication (only required for custom IDP) |
| BTP_TLS_CLIENT_CERTIFICATE | PEM encoded certificate (only required for x509 authentication) |
| BTP_TLS_CLIENT_KEY | PEM encoded private key (only required for x509 authentication) |
| BTP_TLS_IDP_URL | The URL of the identity provider to be used for authentication (only required for x509 authentication) |
| BTP_CLI_SERVER_URL | The URL of the BTP CLI server - **Relevant for SAP internal use-cases only**  |
| BTP_ENABLE_SSO | To use Single Sign-On (SSO) for authentication set this variable to true |

### Terraform provider for Cloud Foundry

For the scenarios where you want to import resources defined in the [Terraform provider for Cloud Foundry](https://registry.terraform.io/providers/cloudfoundry/cloudfoundry/latest/docs) the following parameters are available:

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
