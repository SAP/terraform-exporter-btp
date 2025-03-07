# Resuming Failed Exports

The export of existing infrastructure via the Terraform Exporter for SAP BTP could run into errors during the export process. This could have various reasons. There could be temporal network issues or there might be a temporary issue on the platform.
This leads to an incomplete export of your infrastructure configuration. This is an issue especially if this is an extensive export with many resources that takes some time.

To avoid manual workarounds like restarting the export from scratch or executing a second seperate export with a adjusted configuration to export the missing resources accompnied by some manual rework to get everything in shape, we provide an option to resume the export from the last safepoint.

Let us assume that are exporting a subaccount based on a JSON configuration called `btpResources.json`. The resources that you want to export are:
- a subaccount
- an entitlement
- a subscription
- a service instance

In addition you want to store the generated code into the directory `exported-configuration`. Consequnetly you execute the command:

```bash
btptf export-by-json -s 12345678-abcd-efgh-ae24-86ff3384cf93 -p btpResources.json -c exported-configuration
```

The export process starts, but then runs into an error:

TODO picture

The created files show that the export was interrupted, as the configuration is not complete:

In addition a temporary directory `subscriptions-config` is visible:

TODO Picture

!!! info
    Directories that follow the naing convention `<resource>-config` are part of the export process as a temporary directory contaning the resouce specific information. They get removed when the reosurce was executed successfully or in case of an error that is handled by the Terraform Exporter for SAP BTP. However, in error situations you might see these directories. It is safe to delete them.


You also recognize a file called `importlog.json` in the `exported-configuration` directory:

TODO Picture

The Terraform Exporter for SAP BTP uses this file to track the successfully exported resources.

!!! info
    This file will be removed if the export was successful.

To resume the processing we do not make any changes to the files, but execute the original command:

```bash
btptf export-by-json -s 12345678-abcd-efgh-ae24-86ff3384cf93 -p btpResources.json -c exported-configuration
```

The Terraform Exporter for SAP BTP recognizes the file and prompts how we want to proceed

TODO picture

We select the option to resume the processing.

The processing starts, but as we also have the directory `subscriptions-config`, we get prompted if this should be removed:

It is safe to remove it, so we select the corresponfing option and the export process continues.


The output show that the missing resources get exported:

The summary table gives an overview over all exported resources combining the information from the previous run with the ones from the resumed run:

As a result the export is executed successfully and all resources are available in the genereted configuration.

!!! info
