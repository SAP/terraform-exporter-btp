# Installation

Download the binary from from the [releases ](https://github.com/SAP/terraform-exporter-btp/releases) section of the GitHub repository.

Select the version that you want to use and download the binary that fits your operating system from the `assets` of the release. We recommend using the latest version.

## Local Build

If you contribute to the Terraform Exporter for SAP BTP or you need a fix that has not been added to a released version, you may want to do a local build:

1. Open the Terraform Exporter for SAP BTP from the [GitHub repository](https://github.com/SAP/terraform-exporter-btp) in the VS Code Editor.

2. We have set up a [devcontainer](https://code.visualstudio.com/docs/devcontainers/tutorial), so reopen the repository in the devcontainer.

3. Open a terminal in VS Code and install the binary by running

   ```bash
    make install
    ```
   This will implicitly trigger a build of the source. If you want to build *without* install, execute `make build`.

4. The system will store the binary as `btptf` (`btptf.exe` in case of Windows) in the default binary path of your Go installation `$GOPATH/bin`.

!!! tip
    You find the value of the GOPATH via `go env GOPATH`
