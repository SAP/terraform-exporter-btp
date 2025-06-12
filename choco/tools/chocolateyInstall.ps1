$toolsDir = Split-Path -Parent $MyInvocation.MyCommand.Definition

# Parameters to replace dynamically in workflow
$url = "https://github.com/SAP/terraform-exporter-btp/releases/download/vVERSION/btptf_VERSION_windows_amd64.exe"

$targetPath = Join-Path $toolsDir 'btptf.exe'


$webFileArgs = @{
    PackageName    = $env:ChocolateyPackageName
    FileFullPath   = $targetPath
    Url64bit       = $url
    Checksum64     = 'CHECKSUM'
    ChecksumType64 = 'sha256'
}

# Download the binary at install time
Get-ChocolateyWebFile @webFileArgs

$validExitCodes = @(0)