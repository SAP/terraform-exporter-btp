## btptf completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(btptf completion bash)

To load completions for every new session, execute once:

#### Linux:

	btptf completion bash > /etc/bash_completion.d/btptf

#### macOS:

	btptf completion bash > $(brew --prefix)/etc/bash_completion.d/btptf

You will need to start a new shell for this setup to take effect.


```
btptf completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --verbose   Display verbose output in the console for debugging.
```

### SEE ALSO

* [btptf completion](btptf_completion.md)	 - Generate the autocompletion script for the specified shell

