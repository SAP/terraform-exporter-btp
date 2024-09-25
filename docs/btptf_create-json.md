## btptf create-json

Store the list of resources in a subaccount into a JSON file

### Synopsis

create-json command compiles a list of all resources in a subaccount and store it into a file.

Examples:

btptf create-json --resources=subaccount,entitlements -s <subaccount-id>
btptf create-json --resources=all -s <subaccount-id> -p <file-name.json>

Valid resources are:
- subaccount
- entitlements
- subscriptions
- environment-instances
- trust-configurations
- roles
- role-collections

OR

- all

Mixing "all" with other resources will throw an error.


```
btptf create-json [flags]
```

### Options

```
  -h, --help               help for create-json
  -p, --json-out string    JSON file for list of resources (default "btpResources.json")
  -r, --resources string   comma seperated string for resources (default "all")
```

### Options inherited from parent commands

```
  -d, --debug               Display debugging output in the console.
  -s, --subaccount string   Id of the subaccount
```

### SEE ALSO

* [btptf](btptf.md)	 - Terraform Exporter for SAP BTP

