
# Terraform exporter for SAP BTP

This document lists the syntax and parameters for all three commands.

## btptf

Terraform Exporter for SAP BTP

### Options

```azdeveloper
  -h, --help      help for btptf
      --verbose   Enable verbose output for debugging
  -v, --version   Print the version number of btptf
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
  -d, --directory string    ID of the directory
  -h, --help                help for create-json
  -o, --org string          ID of the Cloud Foundry org
  -p, --path string         Path to JSON file with list of resources (default "btpResources_<account-id>.json")
  -r, --resources string    Comma-separated list of resources to be included (default "all")
  -s, --subaccount string   ID of the subaccount
```

### Options inherited from parent commands

```azdeveloper
      --verbose   Enable verbose output for debugging
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
      --backend-config strings   Backend configuration
  -b, --backend-path string      Path to the Terraform backend configuration file
      --backend-type string      Type of the Terraform backend
  -c, --config-dir string        Directory for the Terraform code (default "generated_configurations_<account-id>")
  -d, --directory string         ID of the directory
  -h, --help                     help for export
  -o, --org string               ID of the Cloud Foundry org
  -r, --resources string         Comma-separated list of resources to be included (default "all")
  -s, --subaccount string        ID of the subaccount
```

### Options inherited from parent commands

```azdeveloper
      --verbose   Enable verbose output for debugging
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
      --backend-config strings   Backend configuration
  -b, --backend-path string      Path to the Terraform backend configuration file
      --backend-type string      Type of the Terraform backend
  -c, --config-dir string        Directory for the Terraform code (default "generated_configurations_<account-id>")
  -d, --directory string         ID of the directory
  -h, --help                     help for export-by-json
  -o, --org string               ID of the Cloud Foundry org
  -p, --path string              Full path to JSON file with list of resources (default "btpResources_<account-id>.json")
  -s, --subaccount string        ID of the subaccount
```

### Options inherited from parent commands

```azdeveloper
      --verbose   Enable verbose output for debugging
```

### See also

* [Back to top](#btptf)

