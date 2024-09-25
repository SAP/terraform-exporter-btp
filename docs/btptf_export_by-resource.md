## btptf export by-resource

export resources of a subaccount

### Synopsis

by-resource command exports the resources of a subaccount as specified.

Examples:

btptf export by-resource --resources=subaccount,entitlements -s <subaccount-id>
btptf export by-resource --resources=all -s <subaccount-id> -p <file-name.json>

Valid resources are:
- subaccount
- entitlements
- subscriptions
- environment-instances
- trust-configurations
- roles
- role-collections
- service-bindings

OR

- all

Mixing "all" with other resources will throw an error.

```
btptf export by-resource [flags]
```

### Options

```
  -h, --help               help for by-resource
  -r, --resources string   comma seperated string for resources (default "all")
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

