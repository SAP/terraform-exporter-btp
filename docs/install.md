# Installation

You have two options to install the btptf CLI:

1. Download the pre-built binary.
2. Local build

The following sections describe the details for the two options.


## Pre-Built Binary

The easiest way to get the binary is to download from the [releases section](https://github.com/SAP/terraform-exporter-btp/releases) of this repository. Select the version that you want to use and download the binary that fits your operating system from the `assets` of the release. We recommend using the latest version.


## Local Build

If you want to build the binary from scratch, follow these steps:

1. Open [this](https://github.com/SAP/terraform-exporter-btp) repository inside VS Code Editor

2. We have setup a [devcontainer](https://code.visualstudio.com/docs/devcontainers/tutorial), so reopen the repository in the devcontainer. We provide devcontainers for Terraform and OpenTofu.

3. Open a terminal in VS Code and install the binary by running

   ```bash
    make install
    ```
   This will implicitly trigger a build of the source. If you want to build *without* install, execute `make build`.

4. The system will store the binary as `btptf` (`btptf.exe` in case of Windows) in the default binary path of your Go installation `$GOPATH/bin`.

!!! tip
    You find the value of the GOPATH via `go env GOPATH`
