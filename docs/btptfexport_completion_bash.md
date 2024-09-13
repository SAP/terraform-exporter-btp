## btptfexport completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(btptfexport completion bash)

To load completions for every new session, execute once:

#### Linux:

	btptfexport completion bash > /etc/bash_completion.d/btptfexport

#### macOS:

	btptfexport completion bash > $(brew --prefix)/etc/bash_completion.d/btptfexport

You will need to start a new shell for this setup to take effect.


```
btptfexport completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [btptfexport completion](btptfexport_completion.md)	 - Generate the autocompletion script for the specified shell

