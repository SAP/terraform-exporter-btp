$toolsDir = Split-Path -Parent $MyInvocation.MyCommand.Definition

# Parameters to replace dynamically in workflow
$url = "https://github.com/SAP/terraform-exporter-btp/releases/download/vVERSION/btptf_VERSION_windows_amd64.exe"

$targetPath = Join-Path $toolsDir 'btptf.exe'

# Download the binary at install time
Get-ChocolateyWebFile `
  -PackageName 'btptf' `
  -FileFullPath $targetPath `
  -Url64bit $url `

# Copy binary to Chocolatey tools location for execution
Copy-Item $targetPath -Destination "$env:ChocolateyToolsLocation\btptf.exe" -Force