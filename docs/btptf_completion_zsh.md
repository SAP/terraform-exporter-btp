## btptf completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(btptf completion zsh)

To load completions for every new session, execute once:

#### Linux:

	btptf completion zsh > "${fpath[1]}/_btptf"

#### macOS:

	btptf completion zsh > $(brew --prefix)/share/zsh/site-functions/_btptf

You will need to start a new shell for this setup to take effect.


```
btptf completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -d, --debug               Display debugging output in the console.
  -s, --subaccount string   Id of the subaccount
```

### SEE ALSO

* [btptf completion](btptf_completion.md)	 - Generate the autocompletion script for the specified shell

