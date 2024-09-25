## btptf export by-json

export resources based on a JSON file.

### Synopsis

Use this command to export resources from the JSON file that is generated using the create-json command.

```
btptf export by-json [flags]
```

### Options

```
  -p, --from string   path to JSON file with resources (default "btpResources.json")
  -h, --help          help for by-json
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

