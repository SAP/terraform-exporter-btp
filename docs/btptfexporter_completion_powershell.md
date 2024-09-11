## btptfexporter completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	btptfexporter completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
btptfexporter completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [btptfexporter completion](btptfexporter_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 11-Sep-2024