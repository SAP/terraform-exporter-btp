## btptfexport resource

Export specific btp resources from a subaccount

### Synopsis


This command is used when you need to export specific resources.
By default, it will generate the <resource_name>_import.tf (import file) and resources.tf (resource file) files.
The resources.tf file can be renamed by using the flag --resourceFileName.
The command will fail if a resource file already exists

```
btptfexport resource [flags]
```

### Options

```
  -h, --help   help for resource
```

### SEE ALSO

* [btptfexport](btptfexport.md)	 - Terraform exporter for BTP
* [btptfexport resource all](btptfexport_resource_all.md)	 - export all resources of a subaccount
* [btptfexport resource entitlements](btptfexport_resource_entitlements.md)	 - export entitlements of a subaccount
* [btptfexport resource environment-instances](btptfexport_resource_environment-instances.md)	 - export environment instance of a subaccount
* [btptfexport resource from-file](btptfexport_resource_from-file.md)	 - export resources from a json file.
* [btptfexport resource subaccount](btptfexport_resource_subaccount.md)	 - export subaccount
* [btptfexport resource subscriptions](btptfexport_resource_subscriptions.md)	 - export subscriptions of a subaccount
* [btptfexport resource trust-configurations](btptfexport_resource_trust-configurations.md)	 - export trust configurations of a subaccount
