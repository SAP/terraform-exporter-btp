## btptfexport resource role-collections

export roles collections of a subaccount

### Synopsis

export role-collections will export all the role collections of the given subaccount and generate resource configuration for it

```
btptfexport resource role-collections [flags]
```

### Options

```
  -o, --config-output-dir string   folder for config generation (default "generated_configurations")
  -h, --help                       help for role-collections
  -f, --resourceFileName string    filename for resource config generation (default "resources.tf")
  -s, --subaccount string          Id of the subaccount
```

### Options inherited from parent commands

```
  -d, --debug   Display debugging output in the console. (default: false)
```

### SEE ALSO

* [btptfexport resource](btptfexport_resource.md)	 - Export specific btp resources from a subaccount

