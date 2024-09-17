## btptfexport completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	btptfexport completion fish | source

To load completions for every new session, execute once:

	btptfexport completion fish > ~/.config/fish/completions/btptfexport.fish

You will need to start a new shell for this setup to take effect.


```
btptfexport completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -d, --debug   Display debugging output in the console. (default: false)
```

### SEE ALSO

* [btptfexport completion](btptfexport_completion.md)	 - Generate the autocompletion script for the specified shell

