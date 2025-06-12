# Installation

## Package Manager
You can easily install the `Terraform Exporter for SAP BTP` using popular installation methods based on your operating system. We’ve made it simple with ready-to-use installers and package managers across platforms. 

!!! Commands
    === "Windows"
        We have released `winget` and `Chocolatey` packages. Please refer to the official [Winget](https://learn.microsoft.com/en-us/windows/package-manager/winget/) and [Chocolatey](https://chocolatey.org/about) documentation websites for details.
        **For Winget based installation**
        ```cmd
        winget install SAP.btptf
        ```
        **For Chocolatey based installation**
        ```cmd
        choco install sap-btptf
        ```

    === "Mac OS"
        Please refer to the [Homebrew](https://brew.sh/) website for details.
        ```bash
        brew tap SAP/terraform-exporter-btp https://github.com/SAP/terraform-exporter-btp
        brew install sap/terraform-exporter-btp/btptf
        ```
    === "Linux"
        We’ve released `deb` and `rpm` packages to support installation on the most common Linux distributions. You can download the packages from the `assets` section of the [releases](https://github.com/SAP/terraform-exporter-btp/releases) page.

        **For Debian-based distributions (like Ubuntu, Linux Mint, etc.):**
        ```bash
        sudo dpkg -i <path-to-download>/terraform-exporter-btp_<latest-version>_linux_amd64.deb
        ```
        **For RPM-based distributions (like Fedora, RHEL, CentOS, openSUSE):**
        ```bash
        sudo rpm -i <path-to-download>/terraform-exporter-btp_<latest-version>_linux_amd64.rpm
        ```

## Download Binaries
You can also download the binaries directly from the [releases](https://github.com/SAP/terraform-exporter-btp/releases) section of the GitHub repository.

Select the version that you want to use and download the binary that fits your operating system from the `assets` of the release. We recommend using the latest version.

## Local Build

If you contribute to the btptf CLI or you need a fix that has not been added to a released version, you may want to do a local build:

1. Open the code of the btptf CLI cloned from the [GitHub repository](https://github.com/SAP/terraform-exporter-btp) in the VS Code Editor.

2. We have set up a [devcontainer](https://code.visualstudio.com/docs/devcontainers/tutorial), so reopen the repository in the devcontainer.

3. Open a terminal in VS Code and install the binary by running `make install`. This will implicitly trigger a build of the source. If you want to build *without* install, execute `make build`.

4. The system will store the binary as `btptf` (`btptf.exe` in case of Windows) in the default binary path of your Go installation `$GOPATH/bin`.

!!! tip
    You find the value of the GOPATH via `go env GOPATH`
