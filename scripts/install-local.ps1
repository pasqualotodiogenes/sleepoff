param(
    [switch]$Build,
    [string]$SourcePath,
    [switch]$Force
)

$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$installDir = Join-Path $env:LOCALAPPDATA "Programs\\sleepoff"
$targetExe = Join-Path $installDir "sleepoff.exe"

if ($Build) {
    Push-Location $repoRoot
    try {
        go build -o sleepoff.exe .
    } finally {
        Pop-Location
    }
    $SourcePath = Join-Path $repoRoot "sleepoff.exe"
}

if (-not $SourcePath) {
    $SourcePath = Join-Path $repoRoot "sleepoff.exe"
}

$resolvedSource = Resolve-Path $SourcePath -ErrorAction SilentlyContinue
if (-not $resolvedSource) {
    throw "Executable not found. Build first with 'go build -o sleepoff.exe .' or run this script with -Build."
}

if ((Test-Path $targetExe) -and (-not $Force)) {
    Write-Host "Replacing existing install at $targetExe"
}

New-Item -ItemType Directory -Path $installDir -Force | Out-Null
Copy-Item -Path $resolvedSource -Destination $targetExe -Force

$envKey = "HKCU:\\Environment"
$currentPath = (Get-ItemProperty -Path $envKey -Name Path -ErrorAction SilentlyContinue).Path
if (-not $currentPath) {
    $currentPath = ""
}

$normalizedInstallDir = $installDir.TrimEnd("\\")
$entries = $currentPath -split ";" | Where-Object { $_.Trim() -ne "" }
$hasInstallDir = $false

foreach ($entry in $entries) {
    if ($entry.Trim().TrimEnd("\\") -ieq $normalizedInstallDir) {
        $hasInstallDir = $true
        break
    }
}

if (-not $hasInstallDir) {
    $newPath = if ($currentPath) { "$currentPath;$installDir" } else { $installDir }
    New-ItemProperty -Path $envKey -Name Path -PropertyType ExpandString -Value $newPath -Force | Out-Null
}

$env:Path = "$installDir;$env:Path"

Add-Type @"
using System;
using System.Runtime.InteropServices;
public static class NativeMethods {
  [DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
  public static extern IntPtr SendMessageTimeout(IntPtr hWnd, uint Msg, IntPtr wParam, string lParam, uint fuFlags, uint uTimeout, out IntPtr lpdwResult);
}
"@

$HWND_BROADCAST = [IntPtr]0xffff
$WM_SETTINGCHANGE = 0x001A
$SMTO_ABORTIFHUNG = 0x0002
$result = [IntPtr]::Zero
[void][NativeMethods]::SendMessageTimeout($HWND_BROADCAST, $WM_SETTINGCHANGE, [IntPtr]::Zero, "Environment", $SMTO_ABORTIFHUNG, 5000, [ref]$result)

Write-Host ""
Write-Host "sleepoff instalado em: $targetExe"
Write-Host "Abra um novo terminal e rode: sleepoff --help"
