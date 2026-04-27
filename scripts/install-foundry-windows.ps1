<#
.SYNOPSIS
  下载 Foundry 官方 Windows 包（win32_amd64.zip）到本机用户目录，并加入当前用户 PATH。不依赖 WSL。

.DESCRIPTION
  在普通 PowerShell 即可运行：
    cd G:\web3\crowdfunding
    powershell -ExecutionPolicy Bypass -File .\scripts\install-foundry-windows.ps1

  若你曾在 PowerShell Profile 里定义了 function forge { bash -lc ... }，请编辑 Profile 删掉或改名该函数，否则可能仍走 WSL。
#>

$ErrorActionPreference = 'Stop'

# PowerShell 5.1 on Windows often does not support Tls13 here.
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$headers = @{
    'User-Agent' = 'Crowdfunding-Foundry-Install'
    'Accept'     = 'application/vnd.github+json'
}

Write-Host 'Fetching latest Foundry release...'
$release = Invoke-RestMethod -Uri 'https://api.github.com/repos/foundry-rs/foundry/releases/latest' -Headers $headers
$zip = $release.assets | Where-Object { $_.name -match 'foundry_.*_win32_amd64\.zip$' } | Select-Object -First 1
if (-not $zip) {
    throw 'Latest release has no foundry_*_win32_amd64.zip. Check https://github.com/foundry-rs/foundry/releases'
}

$destRoot = Join-Path $env:USERPROFILE '.foundry\windows'
New-Item -ItemType Directory -Path $destRoot -Force | Out-Null

$zipPath = Join-Path $env:TEMP $zip.name
Write-Host "Downloading $($zip.name)..."
Invoke-WebRequest -Uri $zip.browser_download_url -OutFile $zipPath -Headers $headers

Write-Host "Extracting to $destRoot ..."
if (Test-Path $destRoot) {
    Get-ChildItem -Path $destRoot -Force -ErrorAction SilentlyContinue | Remove-Item -Recurse -Force -ErrorAction SilentlyContinue
}
New-Item -ItemType Directory -Path $destRoot -Force | Out-Null
Expand-Archive -Path $zipPath -DestinationPath $destRoot -Force

$forge = Get-ChildItem -Path $destRoot -Filter forge.exe -Recurse | Select-Object -First 1
if (-not $forge) {
    throw 'forge.exe not found after extract.'
}
$binDir = $forge.DirectoryName

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ([string]::IsNullOrEmpty($userPath)) {
    $newPath = $binDir
} elseif ($userPath -notlike "*$binDir*") {
    $newPath = $userPath.TrimEnd(';') + ';' + $binDir
} else {
    $newPath = $userPath
}

if ($newPath -ne $userPath) {
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    Write-Host "Added to user PATH: $binDir"
} else {
    Write-Host "User PATH already contains: $binDir"
}

Write-Host ''
Write-Host 'Done. Open a NEW terminal and run:'
Write-Host '  forge --version'
Write-Host '  anvil --version'
Write-Host ''
Write-Host 'If `forge` still runs WSL, remove the `function forge` shim from your PowerShell profile:  notepad $PROFILE'
