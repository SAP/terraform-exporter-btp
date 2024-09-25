## btptfexport resource service-bindings

export service bindings of a subaccount

### Synopsis

export service-bindings will export all the service bindings of the given subaccount and generate resource configuration for it

```
btptfexport resource service-bindings [flags]
```

### Options

```
  -o, --config-output-dir string   folder for config generation (default "generated_configurations")
  -h, --help                       help for service-bindings
  -f, --resourceFileName string    filename for resource config generation (default "resources.tf")
  -s, --subaccount string          Id of the subaccount
```

### Options inherited from parent commands

```
  -d, --debug   Display debugging output in the console. (default: false)
```

### SEE ALSO

* [btptfexport resource](btptfexport_resource.md)	 - Export specific btp resources from a subaccount

