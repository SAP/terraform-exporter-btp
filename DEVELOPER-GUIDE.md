# Developer Guide

TBD

## Debug the CLI

We provide a preconfigured configuration for debugging the btptfexporter commands in VS Code. The configuration is available in the `.vscode` directory as `launch.json`

Here is an example on how to debug the command `btptfexporter resource all`:

1. Set a breakpoint in the file `cmd/exportAll.go` in the run section of the command:

![]()

1. Adjust the `launch.json` configuration to consider your environment variable values. The default are single variables using SSO in the root of the repository:

![]()

> [!WARNING]
> The environment values will be displayed as clear text in the debug console. If you are using your password as environment paramater this will become visible when you start debugging. We therefore highly recommend to use the SSO option.

1. Open the debug pane in VS Code:

![]()

1. Select the configuration `Debug CLI command`:

![]()

1. Run the selection by pressing the green triangle:

![]()

1. VS COde will prompt you for the command via the command palette. It defaults to `resource all -s`:

![]()

1. Enter the command and the parameters you want to use for the command execution. In our case we add a subaccount ID:

![]()

1. Confirm by pressing `Enter`
1. The debugger will start and hit the breakpoint:

![]()

Happy debugging!

## Generate markdown documentation

We can generate the markdown documentation via the make file:

```bash
make docs
```
