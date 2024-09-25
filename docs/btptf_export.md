## btptf export

Export specific resources from an SAP BTP subaccount

### Synopsis


This command is used when you want to export resources of SAP BTP.

You have two options:

- by-json: export resources from a json file that is generated using the create-list command.
- by-resource: export resources you specify by type.

By default, the CLI it will generate the import files and a resource configuration file.
The directory for the configuration files has as default value 'generated_configurations'.
The resource configuration file has as default value 'btp_resources.tf'.

You can change the default values for the directory by using the flag --config-dir.
You can change the name of the resource configuration file by using the flag --resource-file-name.


The command will fail if a resource file already exists

```
btptf export [flags]
```

### Options

```
  -o, --config-dir string           folder for config generation (default "generated_configurations")
  -h, --help                        help for export
  -f, --resource-file-name string   filename for resource config generation (default "btp_resources.tf")
```

### Options inherited from parent commands

```
  -d, --debug               Display debugging output in the console.
  -s, --subaccount string   Id of the subaccount
```

### SEE ALSO

* [btptf](btptf.md)	 - Terraform Exporter for SAP BTP
* [btptf export by-json](btptf_export_by-json.md)	 - export resources based on a JSON file.
* [btptf export by-resource](btptf_export_by-resource.md)	 - export resources of a subaccount

