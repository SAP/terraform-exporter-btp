## btptfexporter resource from-file

export resources from a json file.

### Synopsis

Use this command to export resources from the json file that is generated using the get command.
You can removes resource names from this config file, if you want to selectively import resources

```
btptfexporter resource from-file [flags]
```

### Options

```
  -o, --config-output-dir string    folder for config generation (default "generated_configurations")
  -h, --help                        help for from-file
      --resource-file-path string   json file having subaccount resources info
  -f, --resourceFileName string     filename for resource config generation (default "resources.tf")
  -s, --subaccount string           Id of the subaccount
```

### SEE ALSO

* [btptfexporter resource](btptfexporter_resource.md)	 - Export specific btp resources from a subaccount

###### Auto generated by spf13/cobra on 11-Sep-2024