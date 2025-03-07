# Resuming Failed Exports

The export of existing infrastructure via the Terraform Exporter for SAP BTP could run into errors during the export process. This could have various reasons. There could be temporal network issues or there might be a temporary issue on the platform.
This leads to an incomplete export of your infrastructure configuration. This is an issue especially if this is an extensive export with many resources that takes some time.

To avoid manual workarounds like:
-	restarting the export from scratch.
-	executing a second separate export with an adjusted configuration to export the missing resources accompanied by some manual rework to get everything in shape.

To avoid this we provide an option to resume the export from the last safepoint.

Let us assume that are exporting a subaccount based on a JSON configuration called `btpResources.json`. The resources that we want to export are:
- a subaccount
- an entitlement
- a subscription
- a service instance

In addition, you want to store the generated code into the directory `exported configuration`. Consequently, we execute the command:

```bash
btptf export-by-json -s 12345678-abcd-efgh-ae24-86ff3384cf93 -p btpResources.json -c exported-configuration
```

The export process starts, but then runs into an error:

![Console Output of failed export run](img/export-failed.png)


The created files show that the export was interrupted, as the configuration is not complete:

![Created files after failed export run](img/incomplete-export-files.png)

In addition, a temporary directory `subscriptions-config` is stored.

!!! info
    Directories that follow the naming convention `<resource>-config` are part of the export process as a temporary directory containing the resource specific information. They get removed when the resource was executed successfully or in case of an error that is handled by the Terraform Exporter for SAP BTP. However, in error situations you might see these directories. It is safe to delete them.

You also recognize a file called `importlog.json` in the `exported-configuration` directory:

![File importlog.json as part of failed export run](img/incomplete-export-log.png)

The Terraform Exporter for SAP BTP uses this file to track the successfully exported resources.

!!! info
    This file will be removed if the export was successful.

To resume the processing, we do not make any changes to the files, but execute the original command:

```bash
btptf export-by-json -s 12345678-abcd-efgh-ae24-86ff3384cf93 -p btpResources.json -c exported-configuration
```

The Terraform Exporter for SAP BTP recognizes the file and prompts how we want to proceed

![Repeat export - prompt for resume based on log](img/resume-step1.png)

We select the option to resume the processing.

The processing starts, but as we also have the directory `subscriptions-config`, we get prompted if this should be removed:

![Repeat export - prompt for removing temporary directory](img/resume-step2.png)

It is safe to remove it, so we select the corresponding option, and the export process continues.

The output show that the missing resources get exported:

![Repeat export - console output after success](img/resumed-export-result.png)

The summary table gives an overview over all exported resources combining the information from the previous run with the ones from the resumed run.

As a result, the export is executed successfully, and all resources are available in the generated configuration:

![Repeat export - created files](img/resume-result-files.png)

!!! warning
    Exports can fail due to different reason. While we aim to resume such failed scenarios, there might always be situation that we cannot handle automatically. In this case we advise to restart the export from scratch.
