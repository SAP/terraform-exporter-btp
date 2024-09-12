## btptfexporter generate-resources-list

Store the list of resources in a subaccount into a json file

### Synopsis

generate-resources-list command will get all the resource list or specified resource list in a subaccount.
It will then store this list into a file.

For example:

btptfexporter generate-resources-list --resources=subaccount,entitlements -s <subaccount-id>
btptfexporter generate-resources-list --resources=all -s <subaccount-id> -j <file-name.json>

Valid resources are:
- subaccount
- entitlements
- subscriptions
- environment-instances
- trust-configurations

OR

- all

Mixing "all" with other resources will throw an error.


```
btptfexporter generate-resources-list [flags]
```

### Options

```
  -h, --help                help for generate-resources-list
  -j, --json-out string     json file for list of resources (default "btpResources.json")
  -r, --resources string    comma seperated string for resources (default "all")
  -s, --subaccount string   Id of the subaccount
```

### SEE ALSO

* [btptfexporter](btptfexporter.md)	 - Terraform exporter for BTP

