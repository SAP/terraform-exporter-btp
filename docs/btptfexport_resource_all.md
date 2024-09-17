## btptfexport resource all

export all resources of a subaccount

### Synopsis

export all will export all the resources from a subaccount. Currently only few resources are supported.

export all is a single command to export btp_subaccount, btp_subaccount_entitlements, btp_subaccount_instances, btp_subaccount_subscriptions,
btp_subaccount_trust_configurations 

```
btptfexport resource all [flags]
```

### Options

```
  -o, --config-output-dir string   folder for config generation (default "generated_configurations")
  -h, --help                       help for all
  -f, --resourceFileName string    filename for resource config generation (default "resources.tf")
  -s, --subaccount string          Id of the subaccount
```

### Options inherited from parent commands

```
  -d, --debug   Display debugging output in the console. (default: false)
```

### SEE ALSO

* [btptfexport resource](btptfexport_resource.md)	 - Export specific btp resources from a subaccount

