## btptf completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	btptf completion fish | source

To load completions for every new session, execute once:

	btptf completion fish > ~/.config/fish/completions/btptf.fish

You will need to start a new shell for this setup to take effect.


```
btptf completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -d, --debug               Display debugging output in the console.
  -s, --subaccount string   Id of the subaccount
```

### SEE ALSO

* [btptf completion](btptf_completion.md)	 - Generate the autocompletion script for the specified shell

