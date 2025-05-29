$toolsDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$targetPath = Join-Path $toolsDir 'btptf.exe'

Remove-Item $targetPath -Force -ErrorAction SilentlyContinue