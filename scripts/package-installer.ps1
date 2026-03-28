param(
    [Parameter(Mandatory = $true)]
    [string]$Version,
    [Parameter(Mandatory = $true)]
    [string]$ExecutablePath,
    [string]$OutputDir = "dist"
)

$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$issScript = Join-Path $repoRoot "packaging\\windows\\sleepoff.iss"

if (-not (Test-Path $issScript)) {
    throw "Inno Setup script not found: $issScript"
}

$exePath = Resolve-Path $ExecutablePath
if (-not $exePath) {
    throw "Executable not found: $ExecutablePath"
}

if ([System.IO.Path]::IsPathRooted($OutputDir)) {
    $resolvedOutputDir = $OutputDir
} else {
    $resolvedOutputDir = Join-Path $repoRoot $OutputDir
}

New-Item -ItemType Directory -Path $resolvedOutputDir -Force | Out-Null

$iscc = Get-Command ISCC.exe -ErrorAction SilentlyContinue
if (-not $iscc) {
    $defaultIscc = "C:\\Program Files (x86)\\Inno Setup 6\\ISCC.exe"
    if (Test-Path $defaultIscc) {
        $isccPath = $defaultIscc
    } else {
        throw "ISCC.exe not found. Install Inno Setup 6 and ensure ISCC.exe is available."
    }
} else {
    $isccPath = $iscc.Source
}

$normalizedVersion = $Version.Trim()
if ($normalizedVersion.StartsWith("v")) {
    $normalizedVersion = $normalizedVersion.Substring(1)
}

& $isccPath "/DMyAppVersion=$normalizedVersion" "/DMyAppExe=$exePath" "/DMyOutputDir=$resolvedOutputDir" $issScript
