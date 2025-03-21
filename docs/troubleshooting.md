# Troubleshooting

This page should help you when you run into a situation where the Terraform Exporter for SAP BTP might not produce the expected results.
To narrow down the issue we provide some environment variables that allow you to influence or switch off certain features of the exporter.

The following sections describe the environment variables that you can set to influence the behavior of the exporter depending on the situation you are in.

## The Terraform CLI is used instead of OpenTofu (or vice versa)

You are expecting the the generated code uses the providers from the OpenTofu registry but it downloaded the ones from the HashiCorp Terraform registry or vice versa.

This can happen if you have both CLIs namely the Terraform *and* the OpenTofu CLI installed on your machine. To guide the exporter to use the correct CLI you can set the environment variable `BTPTF_IAC_TOOL` to `terraform` or `tofu` depending on the CLI you want to use.

As an example if you want to use the OpenTofu CLI you can set the environment variable as follows:

=== "Windows"

    ``` powershell
    $env:BTPTF_IAC_TOOL='tofu'
    ```

=== "Linux/Mac"

    ``` bash
    export BTPTF_IAC_TOOL=tofu
    ```

## The generated code looks different than expected

The Terraform Provider for SAP BTP tries to optimize the generated code to make it more readable and maintainable. This includes several actions like removing empty values and introducing implicit and explicit dependencies between resources.

If the generated code looks different than expected you can switch off this feature by setting the environment variable `BTPTF_SKIP_CODECLEANUP` to a non-empty value:

=== "Windows"

    ``` powershell
    $env:BTPTF_SKIP_CODECLEANUP='true'
    ```

=== "Linux/Mac"

    ``` bash
    export BTPTF_SKIP_CODECLEANUP=true
    ```

## Some role collections are not exported

We are removing role collections that are part of the initial creation process of a subaccount or directory or that get created when you subscribe to an application or create a service instance.

If the generated JSON file (command `btptf create-json`) or the generated Terraform code does not contain role collections you would have expected you can deactivate this feature by setting the environment variable `BTPTF_SKIP_RCFILTER` to a non-empty value:

=== "Windows"

    ``` powershell
    $env:BTPTF_SKIP_RCFILTER='true'
    ```

=== "Linux/Mac"

    ``` bash
    export BTPTF_SKIP_RCFILTER=true
    ```

## Some roles are not exported

We are removing roles that are part of the initial creation process of a subaccount or directory or that get created when you subscribe to an application or create a service instance.

If the generated JSON file (command `btptf create-json`) or the generated Terraform code does not contain role collections you would have expected you can deactivate this feature by setting the environment variable `BTPTF_SKIP_ROLEFILTER` to a non-empty value:

=== "Windows"

    ``` powershell
    $env:BTPTF_SKIP_ROLEFILTER='true'
    ```

=== "Linux/Mac"

    ``` bash
    export BTPTF_SKIP_ROLEFILTER=true
    ```
