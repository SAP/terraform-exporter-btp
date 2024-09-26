## btptf export by-json

export resources based on a JSON file.

### Synopsis

Use this command to export resources from the JSON file that is generated using the create-json command.

```
btptf export by-json [flags]
```

### Options

```
  -h, --help          help for by-json
  -p, --path string   path to JSON file with list of resources (default "btpResources.json")
```

### Options inherited from parent commands

```
  -o, --config-dir string           folder for config generation (default "generated_configurations")
  -d, --debug                       Display debugging output in the console.
  -f, --resource-file-name string   filename for resource config generation (default "btp_resources.tf")
  -s, --subaccount string           Id of the subaccount
```

### SEE ALSO

* [btptf export](btptf_export.md)	 - Export specific resources from an SAP BTP subaccount

