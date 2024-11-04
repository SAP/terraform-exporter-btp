
# Terraform exporter for SAP BTP

This document explains the syntax and parameters for the various Terraform exporter for SAP BTP commands.

## btptf

Terraform Exporter for SAP BTP

### Options

```azdeveloper
  -h, --help      help for btptf
      --verbose    Enable verbose output for debugging
```

### See also

* [btptf create-json](#btptf-create-json): Create a JSON file with a list of resources
* [btptf export](#btptf-export): Export resources from SAP BTP
* [btptf export-by-json](#btptf-export-by-json): Export resources from SAP BTP via JSON file

## btptf create-json

Create a JSON file with a list of resources

```bash
btptf create-json [flags]
```

### Options

```azdeveloper
  -d, --directory string      ID of the directory
  -h, --help                  help for create-json
  -o, --organization string   ID of the Cloud Foundry organization
  -p, --path string           Full path to JSON file with list of resources (default "btpResources_<account-id>.json")
  -r, --resources string      Comma-separated list of resources to be included (default "all")
  -s, --subaccount string     ID of the subaccount
```

### Options inherited from parent commands

```azdeveloper
      --verbose    Enable verbose output for debugging
```

### See also

* [Back to top](#btptf)

## btptf export

Export resources from SAP BTP

```bash
btptf export [flags]
```

### Options

```azdeveloper
  -c, --config-dir string     Directory for the Terraform code (default "generated_configurations_<account-id>")
  -d, --directory string      ID of the directory
  -h, --help                  help for export
  -o, --organization string   ID of the Cloud Foundry organization
  -r, --resources string      Comma-separated list of resources to be included (default "all")
  -s, --subaccount string     ID of the subaccount
```

### Options inherited from parent commands

```azdeveloper
      --verbose    Enable verbose output for debugging
```

### See also

* [Back to top](#btptf)

## btptf export-by-json

Export resources from SAP BTP via JSON file

```bash
btptf export-by-json [flags]
```

### Options

```azdeveloper
  -c, --config-dir string     Directory for the Terraform code (default "generated_configurations_<account-id>")
  -d, --directory string      ID of the directory
  -h, --help                  help for export-by-json
  -o, --organization string   ID of the Cloud Foundry organization
  -p, --path string           Full path to JSON file with list of resources (default "btpResources_<account-id>.json")
  -s, --subaccount string     ID of the subaccount
```

### Options inherited from parent commands

```azdeveloper
      --verbose    Enable verbose output for debugging
```

### See also

* [Back to top](#btptf)

