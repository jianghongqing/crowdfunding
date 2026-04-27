<#
.SYNOPSIS
  当 wsl --update / wsl --install 报 0x80072f7d（安全频道）时，从 GitHub 下载官方 WSL x64 MSI 并安装，绕过微软在线更新通道。

.DESCRIPTION
  需要「以管理员身份运行」PowerShell，例如：
    cd G:\web3\crowdfunding
    powershell -ExecutionPolicy Bypass -File .\scripts\fix-wsl.ps1

  可选：若浏览器能上但 WinHTTP 未同步代理，可加 -SyncProxyFromIE；若代理配错可加 -ResetWinHTTPProxy。

.NOTES
  安装完成后请用 Microsoft Store 安装 Ubuntu，或再执行 wsl --install -d Ubuntu（若仍失败，以 Store 为准）。
#>
#Requires -RunAsAdministrator

param(
    [switch]$SyncProxyFromIE,
    [switch]$ResetWinHTTPProxy
)

$ErrorActionPreference = 'Stop'

[Net.ServicePointManager]::SecurityProtocol =
    [Net.SecurityProtocolType]::Tls12 -bor [Net.SecurityProtocolType]::Tls13

if ($ResetWinHTTPProxy) {
    Write-Host 'Resetting WinHTTP proxy...'
    netsh winhttp reset proxy | Out-Null
}
if ($SyncProxyFromIE) {
    Write-Host 'Importing WinHTTP proxy from IE/系统设置...'
    netsh winhttp import proxy source=ie | Out-Null
}

$headers = @{
    'User-Agent' = 'Crowdfunding-WSL-Fix'
    'Accept'     = 'application/vnd.github+json'
}

Write-Host 'Fetching latest WSL release from GitHub...'
$release = Invoke-RestMethod -Uri 'https://api.github.com/repos/microsoft/WSL/releases/latest' -Headers $headers
$msi = $release.assets | Where-Object { $_.name -match '^wsl\..*\.x64\.msi$' } | Select-Object -First 1
if (-not $msi) {
    throw 'Latest WSL release has no wsl.*.x64.msi asset. Open https://github.com/microsoft/WSL/releases manually.'
}

$out = Join-Path $env:TEMP $msi.name
$mb = [math]::Round($msi.size / 1MB, 1)
Write-Host "Downloading $($msi.name) (${mb} MB)..."
Invoke-WebRequest -Uri $msi.browser_download_url -OutFile $out -Headers $headers

Write-Host 'Installing MSI (passive UI)...'
$proc = Start-Process -FilePath 'msiexec.exe' -ArgumentList @('/i', $out, '/passive', '/norestart') -Wait -PassThru
if ($proc.ExitCode -ne 0) {
    throw "msiexec failed with exit code $($proc.ExitCode)"
}

Write-Host 'Shutting down WSL...'
wsl.exe --shutdown 2>$null

Write-Host ''
Write-Host 'wsl --version:'
& wsl.exe --version

Write-Host ''
Write-Host 'Next: Install Ubuntu from Microsoft Store (search Ubuntu 22.04), then open Ubuntu once to finish setup.'
Write-Host 'If Store is OK but CLI still fails: https://aka.ms/wslubuntu2204'
