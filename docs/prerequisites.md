# Prerequisite

After executing the setup of the btptf CLI, you must set some required environment variables needed for authentication.

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
