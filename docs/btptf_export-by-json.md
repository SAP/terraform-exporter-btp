## btptf export-by-json

Export resources from SAP BTP via JSON file

```
btptf export-by-json [flags]
```

### Options

```
  -c, --config-dir string   Directory for the Terraform code (default "generated_configurations_<account-id>")
  -d, --directory string    ID of the directory
  -h, --help                help for export-by-json
  -p, --path string         Full path to JSON file with list of resources (default "btpResources_<account-id>.json")
  -s, --subaccount string   ID of the subaccount
```

### Options inherited from parent commands

```
      --verbose    Enable verbose output for debugging
```

### SEE ALSO

* [btptf](btptf.md)	 - Terraform Exporter for SAP BTP

