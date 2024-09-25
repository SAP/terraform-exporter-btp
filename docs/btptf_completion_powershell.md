## btptf completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	btptf completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
btptf completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -d, --debug               Display debugging output in the console.
  -s, --subaccount string   Id of the subaccount
```

### SEE ALSO

* [btptf completion](btptf_completion.md)	 - Generate the autocompletion script for the specified shell

