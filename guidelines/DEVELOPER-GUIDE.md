# Developer Guide

## Debugging Output

By default the CLI suppresses the verbose output of the different Terraform commands. However, they might be quite useful, when it comes to analyzing issues. For that you can add the `--verbose` flag to any command of the CLI, which will result in the full output of any `cmd.exec()` execution suppressing any "UI candy" like spinners.

## Debug the CLI

We provide a configuration for debugging the `btptf` commands in VS Code. The configuration is available in the `.vscode` directory as `launch.json`

Here is an example of how to debug the command `btptf create-json`:

1. Set a breakpoint in the file `cmd/createJson.go` in the run section of the command:

   <img src="../assets/devguide-pics/debug0.png" alt="Set a breakpoint in VS Code" width="600px">

1. Adjust the `launch.json` configuration to consider your environment variable values. The default are single variables using SSO in the root of the repository:

   <img src="../assets/devguide-pics/debug0b.png" alt="VS Code Debug launch configuration" width="600px">

   > [!WARNING]
   > The environment values will be displayed as clear text in the debug console. If you are using your password as an environment parameter this will become visible when you start debugging. We therefore highly recommend using the SSO option.

1. Open the debug perspective in the VS Code side bar:

   <img src="../assets/devguide-pics/debug1.png" alt="VS Code Side bar" width="50px">

1. Select the configuration `Debug CLI command`:

   <img src="../assets/devguide-pics/debug2.png" alt="VS Code debug configuration options" width="600px">

1. Run the selection by pressing the green triangle:

   <img src="../assets/devguide-pics/debug3.png" alt="Run debug configuration" width="600px">

1. VS Code will prompt you for the command via the command palette. It defaults to `resource all -s`. Enter the command and the parameters that you want to use for the command execution. In our case we add a subaccount ID and confirm by pressing `Enter`:

   <img src="../assets/devguide-pics/debug4.png" alt="Prompt for parameters in debug configuration" width="600px">

1. The debugger will start and hit the breakpoint:

   <img src="../assets/devguide-pics/debug5.png" alt="VS Code hitting breakpoint" width="600px">

Happy debugging!

## Generate Markdown Documentation

When updating command descriptions, you must generate the markdown documentation via the make file:

```bash
make docs
```

## Adding support for New Resources on Subaccount Level

To enable new resources on the subaccount level, you must execute the following steps:

1. Add the corresponding *constants* for the command parameter and the technical resource name in the `tfutils/tfutils.go` file.
1. Add the *mapping of the constants* into the function `TranslateResourceParamToTechnicalName` in the `tfutils/tfutils.go` file.
1. Add the command constant to the slice of `AllowedResources` in the `tfutil/tfconfig.go` file.
1. Create a new implementation for the import factory in the directory `tfimportprovider`. You can take the file `subaccountRoleCollectionImportProvider.go` as an example concerning the structure of the file.
1. Add the new implementation to the import factory function `GetImportBlockProvider` in the file `tfimportprovider/tfImportProviderFactory.go`.
1. Depending on the resource, you must define a transformation of the data from the data source to a string array. Place this logic into the function `transformDataToStringArray` in the `tfutils/tfutils.go` file.
1. Depending on your resource you might also need to add custom formatting logic for the resource address in the Terraform configuration. Place that into the file `output/format.go`. In most cases the function `FormatResourceNameGeneric` is sufficient.

### Adding Unit Tests

The main domain logic that we must test is located in the factory implementations in the directory `tfimportprovider`. Creating these tests should reflect the real-world setup, so we need to extract the test data from subaccounts and store them in the tests. In the following sections, we describe how to best extract this data namely the JSON string that you need as input for your test.

#### Prerequisites

As a prerequisite you should have a Terraform account with the resource that you want to cover in your test up and running. We will use Terraform to extract the base data.

#### Extracting the Data

First create a Terraform setup that allows you to read the data via a data source. The basic setup could look like this:

- `provider.tf`:

```terraform
terraform {

  required_providers {
    btp = {
      source  = "SAP/btp"
      version = "~>1.9.0"
    }
  }
}

provider "btp" {
  globalaccount  = "<YOU GLOBALACCOUNT SUBDOMAIN>"
}
```

- Assuming we want to fetch subscriptions the `main.tf`would look like this

```terraform
data "btp_subaccount_subscriptions" "all" {
  subaccount_id = "<YOUR SUBACCOUNT ID>"
}

output "all" {
  value = data.btp_subaccount_subscriptions.all
}

```

Next we execute a planning and store the plan file:

```bash
terraform plan -out plan.out
```

you have two options now:

- If you want to create a JSON string with all the resources contained in the `plan.out`, execute the script `guidelines/scripts/transform_all.sh` that needs to be located at the same level as the `plan.out` file.
- If you want to adjust the result you must execute the following steps:
   1. Generate the JSON file via: Terraform show -json plan.out |  jq .planned_values.outputs.all.value > restrictedplan.json
   1. Adjust the JSON file e.g., remove some entries
   1. execute the script `guidelines/scripts/transform_json.sh` that needs to be located at the same level as the `restrictedplan.json` file.

With that you get a file that contains the JSON string that you can use as input for your tests of the creation of the import block functions.

#### Creating the Unit Test

An example how to create test case is given by the unit test implemented in `tfimportprovider/subaccountSubscriptionImportProvider_test.go`.
GitHub Copilot can be quite useful to setup the basics for the test, but some rework is needed.

## Creating Console Help

We use the custom templating option available in the Cobra Framework to construct the output in the console. The override of the default templating flow is triggered in the commands via the function `SetHelpTemplate` and `SetUsageTemplate`.

In general, we call the `generateCmdHelp` function to generate the output that will be displayed in the console. The `generateCmdHelp` function gets the command as well as a structure of the type `generateCmdHelpOptions`.

If the command receives an empty structure, it will call several default functions to create the console help. However, you have the option to override the single section by providing a custom function that crafts the string used in the console help.

You find an example for this setup in the command `exportByResourceCmd`. Be aware that the code leverages several helper functions that are available in the file `cmdDocsHelper.go`.
