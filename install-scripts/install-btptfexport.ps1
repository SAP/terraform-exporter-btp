#!/usr/bin/env pwsh
<#
.SYNOPSIS
Download and install btptfexport on the local machine.

.DESCRIPTION
Downloads and installs btptfexport on the local machine. Includes ability to configure
download and install locations.

.PARAMETER BaseUrl
Specifies the base URL to use when downloading. Default is
https://raw.githubusercontent.com/SAP/terraform-exporter-btp/refs/heads/main/install-scripts/install-btptfexport.ps1

.PARAMETER Version
Specifies the version to use. Default is `latest`. Valid values include a
SemVer version number (e.g. 1.0.0 or 0.1.0-beta1), `latest`

.PARAMETER DryRun
Print the download URL and quit. Does not download or install.

.PARAMETER InstallFolder
Location to install btptfexport.

.PARAMETER SymlinkFolder
(Mac/Linux only) Folder to symlink

.PARAMETER DownloadTimeoutSeconds
Download timeout in seconds. Default is 120 (2 minutes).

.PARAMETER InstallShScriptUrl
(Mac/Linux only) URL to the install-btptfexport.ps1 script. Default is https://raw.githubusercontent.com/SAP/terraform-exporter-btp/refs/heads/main/install-scripts/install-btptfexport.ps1

.EXAMPLE
powershell -ex AllSigned -c "Invoke-RestMethod 'https://raw.githubusercontent.com/SAP/terraform-exporter-btp/refs/heads/main/install-scripts/install-btptfexport.ps1' | Invoke-Expression"

Install the btptfexport CLI from a Windows shell

The use of `-ex AllSigned` is intended to handle the scenario where a machine's
default execution policy is restricted such that modules used by
`install-btptfexport.ps1` cannot be loaded. Because this syntax is piping output from
`Invoke-RestMethod` to `Invoke-Expression` there is no direct valication of the
`install-btptfexport.ps1` script's signature. Validation of the script can be
accomplished by downloading the script to a file and executing the script file.

.EXAMPLE
Invoke-RestMethod 'https://raw.githubusercontent.com/SAP/terraform-exporter-btp/refs/heads/main/install-scripts/install-btptfexport.ps1' -OutFile 'install-btptfexport.ps1'
PS > ./install-btptfexport.ps1

Download the installer and execute from PowerShell
#>

param(
    [string] $BaseUrl = "https://btptfexportrelease.azureedge.net/btptfexport/standalone/release",
    [string] $Version = "latest",
    [switch] $DryRun,
    [string] $InstallFolder,
    [string] $SymlinkFolder,
    [switch] $SkipVerify = $true,
    [int] $DownloadTimeoutSeconds = 120,
    [string] $InstallShScriptUrl = "https://raw.githubusercontent.com/SAP/terraform-exporter-btp/refs/heads/main/install-scripts/install-btptfexport.ps1"
)

function isLinuxOrMac {
    return $IsLinux -or $IsMacOS
}

# Does some very basic parsing of /etc/os-release to output the value present in
# the file. Since only lines that start with '#' are to be treated as comments
# according to `man os-release` there is no additional parsing of comments
# Options like:
# bash -c "set -o allexport; source /etc/os-release;set +o allexport; echo $VERSION_ID"
# were considered but it's possible that bash is not installed on the system and
# these commands would not be available.
function getOsReleaseValue($key) {
    $value = $null
    foreach ($line in Get-Content '/etc/os-release') {
        if ($line -like "$key=*") {
            # 'ID="value" -> @('ID', '"value"')
            $splitLine = $line.Split('=', 2)

            # Remove surrounding whitespaces and quotes
            # ` "value" ` -> `value`
            # `'value'` -> `value`
            $value = $splitLine[1].Trim().Trim(@("`"", "'"))
        }
    }
    return $value
}

function getOs {
    $os = [Environment]::OSVersion.Platform.ToString()
    try {
        if (isLinuxOrMac) {
            if ($IsLinux) {
                $os = getOsReleaseValue 'ID'
            } elseif ($IsMacOs) {
                $os = sw_vers -productName
            }
        }
    } catch {
        Write-Error "Error getting OS name $_"
        $os = "error"
    }
    return $os
}

function getOsVersion {
    $version = [Environment]::OSVersion.Version.ToString()
    try {
        if (isLinuxOrMac) {
            if ($IsLinux) {
                $version = getOsReleaseValue 'VERSION_ID'
            } elseif ($IsMacOS) {
                $version = sw_vers -productVersion
            }
        }
    } catch {
        Write-Error "Error getting OS version $_"
        $version = "error"
    }
    return $version
}

function isWsl {
    $isWsl = $false
    if ($IsLinux) {
        $kernelRelease = uname --kernel-release
        if ($kernelRelease -like '*wsl*') {
            $isWsl = $true
        }
    }
    return $isWsl
}

function getTerminal {
    return (Get-Process -Id $PID).ProcessName
}

if (isLinuxOrMac) {
    if (!(Get-Command curl)) {
        Write-Error "Command could not be found: curl."
        exit 1
    }
    if (!(Get-Command bash)) {
        Write-Error "Command could not be found: bash."
        exit 1
    }

    $params = @(
        '--base-url', "'$BaseUrl'",
        '--version', "'$Version'"
    )

    if ($InstallFolder) {
        $params += '--install-folder', "'$InstallFolder'"
    }

    if ($SymlinkFolder) {
        $params += '--symlink-folder', "'$SymlinkFolder'"
    }

    if ($DryRun) {
        $params += '--dry-run'
    }

    if ($VerbosePreference -eq 'Continue') {
        $params += '--verbose'
    }

    $bashParameters = $params -join ' '
    Write-Verbose "Running: curl -fsSL $InstallShScriptUrl | bash -s -- $bashParameters" -Verbose:$Verbose
    bash -c "curl -fsSL $InstallShScriptUrl | bash -s -- $bashParameters"
    exit $LASTEXITCODE
}

try {
    $packageFilename = "btptfexport-windows-amd64.msi"

    $downloadUrl = "$BaseUrl/$packageFilename"
    if ($Version) {
        $downloadUrl = "$BaseUrl/$Version/$packageFilename"
    }

    if ($DryRun) {
        Write-Host $downloadUrl
        exit 0
    }

    $tempFolder = "$([System.IO.Path]::GetTempPath())$([System.IO.Path]::GetRandomFileName())"
    Write-Verbose "Creating temporary folder for downloading executable: $tempFolder"
    New-Item -ItemType Directory -Path $tempFolder | Out-Null

    Write-Verbose "Downloading executable from $downloadUrl" -Verbose:$Verbose
    $releaseArtifactFilename = Join-Path $tempFolder $packageFilename
    try {
        $global:LASTEXITCODE = 0
        Invoke-WebRequest -Uri $downloadUrl -OutFile $releaseArtifactFilename -TimeoutSec $DownloadTimeoutSeconds
        if ($LASTEXITCODE) {
            throw "Invoke-WebRequest failed with nonzero exit code: $LASTEXITCODE"
        }
    } catch {
        Write-Error -ErrorRecord $_
        exit 1
    }

    try {
        Write-Verbose "Verifying signature of $releaseArtifactFilename" -Verbose:$Verbose
        $signature = Get-AuthenticodeSignature $releaseArtifactFilename
        if ($signature.Status -ne 'Valid') {
            Write-Error "Signature of $releaseArtifactFilename is not valid"
            reportTelemetryIfEnabled 'InstallFailed' 'SignatureVerificationFailed'
            exit 1
        }
    } catch {
        Write-Error -ErrorRecord $_
        reportTelemetryIfEnabled 'InstallFailed' 'SignatureVerificationFailed'
        exit 1
    }

    try {
        Write-Verbose "Installing btptfexport CLI" -Verbose:$Verbose
        $MSIEXEC = "${env:SystemRoot}\System32\msiexec.exe"
        $installProcess = Start-Process $MSIEXEC `
            -ArgumentList @("/i", "`"$releaseArtifactFilename`"", "/qn", "INSTALLDIR=`"$InstallFolder`"", "INSTALLEDBY=`"install-btptfexport.ps1`"") `
            -PassThru `
            -Wait

        if ($installProcess.ExitCode) {
            if ($installProcess.ExitCode -eq 1603) {
                Write-Host "A later version of Azure Developer CLI may already be installed. Use 'Add or remove programs' to uninstall that version and try again."
            }

            Write-Error "Could not install MSI at $releaseArtifactFilename. msiexec.exe returned exit code: $($installProcess.ExitCode)"
            exit 1
        }
    } catch {
        Write-Error -ErrorRecord $_
        exit 1
    }

    Write-Verbose "Cleaning temporary install directory: $tempFolder" -Verbose:$Verbose
    Remove-Item $tempFolder -Recurse -Force | Out-Null

    if (!(isLinuxOrMac)) {
        # Installed on Windows
        Write-Host "Successfully installed btptfexport"
        Write-Host "The Terraform Exporter for SAP BTP CLI installed successfully. You may need to restart running programs for installation to take effect."
        Write-Host "- For Windows Terminal, start a new Windows Terminal instance."
        Write-Host "- For VSCode, close all instances of VSCode and then restart it."
    }
    exit 0
} catch {
    Write-Error -ErrorRecord $_
    exit 1
}
